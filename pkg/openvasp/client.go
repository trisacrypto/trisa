package openvasp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client is used to make HTTPS requests with mTLS configured to TRP servers.
type Client struct {
	http.Client
}

func NewClient() *Client {
	cookies, _ := cookiejar.New(nil)
	return &Client{
		Client: http.Client{
			// TODO: update transport with mTLS credentials
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           cookies,
			Timeout:       1 * time.Minute,
		},
	}
}

// Inquiry posts a Travel Rule Inquiry to the counterparty as specified by the info in
// the inquiry request and returns the response from the client.
func (c *Client) Inquiry(inquiry *Inquiry) (_ *TravelRuleResponse, err error) {
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(inquiry); err != nil {
		return nil, err
	}
	return c.Post(inquiry.TRP, body)
}

// Confirm posts a Travel Rule Confirmation to the counterparty as specified by the info
// in the confirmation request and returns the response from the client.
func (c *Client) Confirmation(confirm *Confirmation) (_ *TravelRuleResponse, err error) {
	body := bytes.NewBuffer(nil)
	if err = json.NewEncoder(body).Encode(confirm); err != nil {
		return nil, err
	}
	return c.Post(confirm.TRP, body)
}

// Post a request to the specified Address, which can be a Travel Address, LNURL, or
// plain-text URL. The content type is automatically set to application/json and the
// body should read JSON data. Required OpenVASP headers are set by the TRPInfo struct.
// If the APIVersion is omitted, the default APIVersion is used, similarly if no request
// identifier is present, one is generated for the reuqest.
func (c *Client) Post(info *TRPInfo, body io.Reader) (_ *TravelRuleResponse, err error) {
	var req *http.Request
	if req, err = NewRequest(info, body); err != nil {
		return nil, err
	}

	var rep *http.Response
	if rep, err = c.Client.Do(req); err != nil {
		return nil, err
	}

	trr := &TravelRuleResponse{Response: *rep}
	if serr := trr.StatusError(); serr != nil {
		return nil, serr
	}
	return trr, nil
}

// NewRequest generates a new TRP POST request from the specified info, using the
// Address to determine the endpoint and setting headers accordingly. If the APIVersion
// is omitted, the default APIVersion is used. If there is no request identifier, one
// is generated and added to the info struct.
func NewRequest(info *TRPInfo, body io.Reader) (req *http.Request, err error) {
	if info.APIVersion == "" {
		info.APIVersion = APIVersion
	}

	if info.RequestIdentifier == "" {
		info.RequestIdentifier = uuid.NewString()
	}

	var endpoint string
	if endpoint, err = info.GetURL(); err != nil {
		return nil, err
	}

	if req, err = http.NewRequest(http.MethodPost, endpoint, body); err != nil {
		return nil, err
	}

	req.Header.Set(ContentTypeHeader, ContentTypeValue)
	req.Header.Set(APIVersionHeader, info.APIVersion)
	req.Header.Set(RequestIdentifierHeader, info.RequestIdentifier)
	if len(info.APIExtensions) > 0 {
		req.Header.Set(APIExtensionsHeader, strings.Join(info.APIExtensions, ","))
	}
	return req, nil
}

// TravelRuleResponse embeds an http Response and makes provisions for parsing the
// response data back from the counterparty.
type TravelRuleResponse struct {
	http.Response
	info        *TRPInfo
	statusError *StatusError
}

// Info returns the travel rule info from the headers of the response
func (t *TravelRuleResponse) Info() *TRPInfo {
	if t.info == nil {
		t.info = &TRPInfo{
			Address:           t.Request.URL.String(),
			APIVersion:        t.Header.Get(APIVersionHeader),
			RequestIdentifier: t.Header.Get(RequestIdentifierHeader),
		}

		if extensions := t.Header.Get(APIExtensionsHeader); extensions != "" {
			t.info.APIExtensions = strings.Split(extensions, ",")
		}
	}
	return t.info
}

// If a 400 or 500 status code was received, the status error is returned.
func (t *TravelRuleResponse) StatusError() *StatusError {
	if t.statusError == nil {
		if t.StatusCode < 200 || t.StatusCode >= 300 {
			t.statusError = &StatusError{
				Code: t.StatusCode,
			}

			contentType := t.Header.Get(ContentTypeHeader)
			if contentType != "" {
				if mediaType, _, err := mime.ParseMediaType(contentType); err == nil {
					if mediaType == "text/plain" {
						msg, _ := io.ReadAll(t.Body)
						t.statusError.Message = strings.TrimSpace(string(msg))
						t.Body.Close()
					}
				}
			}

		}
	}
	return t.statusError
}

// InquiryResolution attempts to parse a travel rule response from the response body,
// closing the response body when it is complete. If the body has already been closed
// then an EOF error will be returned.
func (t *TravelRuleResponse) InquiryResolution() (*InquiryResolution, error) {
	// Check the status and the header before reading the body
	// There is a possibility that nil, nil is returned if the server returned a 204.
	if t.StatusCode != http.StatusOK {
		return nil, t.StatusError()
	}

	// Check the content type to ensure JSON data was returned
	contentType := t.Header.Get(ContentTypeHeader)
	if contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}

		if mediaType != ContentMediaType {
			return nil, fmt.Errorf("could not parse response content type %s", contentType)
		}
	}

	defer t.Body.Close()

	// Create the JSON decoder
	decoder := json.NewDecoder(t.Body)
	decoder.DisallowUnknownFields()

	resolution := &InquiryResolution{}
	if err := decoder.Decode(resolution); err != nil {
		return nil, err
	}

	return resolution, nil
}
