package openvasp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

const (
	// Maximum size in bytes for a travel rule payload: 10MiB
	MaxPayloadSize = 1.049e+7
)

func TransferInquiry(handler InquiryHandler) http.Handler {
	return APIChecks(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the travel rule inquiry
		var inquiry *trp.Inquiry
		if err := decodeJSON(w, r, &inquiry); err != nil {
			httpError(w, err)
			return
		}

		// Add the TRP Info to the inquiry from the headers
		inquiry.Info = ParseTRPInfo(r)

		// TODO: validate the inquiry received

		out, err := handler.OnInquiry(inquiry)
		if err != nil {
			httpError(w, err)
			return
		}

		// If out is nil and no error is specified send standard response.
		if out == nil {
			out = &trp.Resolution{}
		}

		// If not automatically approved or rejected, add the API version to the reply.
		if out.Approved == nil && out.Rejected == "" {
			out.Version = APIVersion
		}

		// Default response is 200 with the API Version
		w.Header().Set(ContentTypeHeader, ContentTypeValue)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(out)
	}))
}

func TransferConfirmation(handler ConfirmationHandler) http.Handler {
	return APIChecks(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode the confirmation message
		var confirmation trp.Confirmation
		if err := decodeJSON(w, r, &confirmation); err != nil {
			httpError(w, err)
			return
		}

		// Add the TRP Info to the confirmation from the headers
		confirmation.Info = ParseTRPInfo(r)

		// Validate the confirmation message
		if err := confirmation.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := handler.OnConfirmation(&confirmation); err != nil {
			httpError(w, err)
			return
		}

		// If the confirmation is successful then a 204 no-content is returned.
		w.WriteHeader(http.StatusNoContent)
	}))
}

// APIChecks is middleware that asserts that the headers in the TRP request are correct
// and valid, ensuring that the core protocol is implemented correctly.
func APIChecks(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A POST request is expected both for inquiries and confirmations.
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// Enforce Application Version
		apiVersion := r.Header.Get(APIVersionHeader)
		if apiVersion != APIVersion {
			http.Error(w, "must specify api version header "+APIVersion, http.StatusBadRequest)
			return
		}

		// Set the APIVersion header in the outgoing response
		w.Header().Add(APIVersionHeader, APIVersion)

		// Must specify a request identifier
		var requestIdentifier string
		if requestIdentifier = r.Header.Get(RequestIdentifierHeader); requestIdentifier == "" {
			http.Error(w, "must specify request identifier", http.StatusBadRequest)
			return
		}

		// Echo back the request identifier in the outgoing response
		w.Header().Add(RequestIdentifierHeader, requestIdentifier)

		// Enforce JSON content type; if no content-type is specified assume JSON
		contentType := r.Header.Get(ContentTypeHeader)
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "malformed content-type header", http.StatusUnsupportedMediaType)
				return
			}

			if mt != MIMEJSON {
				http.Error(w, "content-type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Parse TRPInfo from the headers of an HTTP request. If any headers are not present,
// then the info is populated with assumed fields or empty values as appropriate.
// TODO: parse the LNURL from the URL rather than passing the raw URL.
func ParseTRPInfo(r *http.Request) *trp.Info {
	info := &trp.Info{
		Address:           r.URL.String(),
		APIVersion:        r.Header.Get(APIVersionHeader),
		RequestIdentifier: r.Header.Get(RequestIdentifierHeader),
	}

	// TODO: do we need escaping or more extensive parsing?
	if extensions := r.Header.Get(APIExtensionsHeader); extensions != "" {
		info.APIExtensions = strings.Split(extensions, ",")
	}

	return info
}

func httpError(w http.ResponseWriter, err error) {
	var status *trp.StatusError
	if errors.As(err, &status) {
		http.Error(w, status.Error(), status.Code)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Convert the body into a max bytes reader to limit JSON payloads for security
	r.Body = http.MaxBytesReader(w, r.Body, MaxPayloadSize)

	// Create the JSON decoder
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the JSON data and handle any errors
	if err := decoder.Decode(&dst); err != nil {
		var (
			syntaxError *json.SyntaxError
			typeError   *json.UnmarshalTypeError
			maxBytes    *http.MaxBytesError
		)

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &trp.StatusError{Code: http.StatusBadRequest, Message: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &trp.StatusError{Code: http.StatusBadRequest, Message: "request body contains badly-formed JSON"}

		case errors.As(err, &typeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at %d)", typeError.Field, typeError.Offset)
			return &trp.StatusError{Code: http.StatusBadRequest, Message: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("request body contains unknown field %s", fieldName)
			return &trp.StatusError{Code: http.StatusBadRequest, Message: msg}

		case errors.Is(err, io.EOF):
			return &trp.StatusError{Code: http.StatusBadRequest}

		case errors.As(err, &maxBytes):
			return &trp.StatusError{Code: http.StatusRequestEntityTooLarge}

		default:
			return err
		}
	}

	// Ensure the request body only contains a single JSON object
	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return &trp.StatusError{Code: http.StatusBadRequest, Message: "request body must only contain a single JSON object"}
	}

	return nil
}
