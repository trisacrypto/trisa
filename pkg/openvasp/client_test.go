package openvasp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/lnurl"
)

func TestClientPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the method
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// Check the headers
		if vers := r.Header.Get(openvasp.APIVersionHeader); vers != openvasp.APIVersion {
			http.Error(w, "unexpected version", http.StatusBadRequest)
			return
		}

		if ct := r.Header.Get(openvasp.ContentTypeHeader); ct != openvasp.ContentTypeValue {
			http.Error(w, "unexpected content type", http.StatusUnsupportedMediaType)
			return
		}

		if id := r.Header.Get(openvasp.RequestIdentifierHeader); id == "" {
			http.Error(w, "no request identifier", http.StatusBadRequest)
			return
		}

		if ae := r.Header.Get(openvasp.APIExtensionsHeader); ae == "" {
			http.Error(w, "no api extensions header", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	// TODO: cannot test travel rule addresses because a) they require TLS and b) they
	// do not support having ports in the URL at the moment.
	client := openvasp.NewClient()
	info := &openvasp.TRPInfo{
		Address:       srv.URL + "/x/12345?t=i",
		APIExtensions: []string{openvasp.UnsealedTRISAExtension},
	}

	req, err := client.Post(info, nil)
	require.NoError(t, err, "could not execute post to normal url")
	require.Equal(t, http.StatusNoContent, req.StatusCode, "expected a 200 with a regular url")

	// Use a LNURL to make the request
	info.Address, err = lnurl.Encode(info.Address)
	require.NoError(t, err, "could not lnurl encode the address")

	req, err = client.Post(info, nil)
	require.NoError(t, err, "could not execute post to normal url")
	require.Equal(t, http.StatusNoContent, req.StatusCode, "expected a 200 with a regular url")
}

func TestClientInquiry(t *testing.T) {
	// Setup the mock handler the http test server.
	mock := &MockHandler{}
	handler := openvasp.TransferInquiry(mock)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// Set up the test fixtures and client
	inquiry, err := loadInquiryPayload("testdata/inquiry.json")
	require.NoError(t, err, "could not load fixture from testdata/inquiry.json")

	client := openvasp.NewClient()
	inquiry.TRP = &openvasp.TRPInfo{
		Address:           srv.URL + "/x/12345?t=i",
		RequestIdentifier: "ba3bc19e-50f8-4e48-8a2a-da1e50f9b672",
	}

	t.Run("Error", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallInquiry, errors.New("compliance overload"))

		rep, err := client.Inquiry(inquiry)
		require.Nil(t, rep, "response should be nil in the case of an error")
		require.Error(t, err, "expected an error to be returned")

		serr, ok := err.(*openvasp.StatusError)
		require.True(t, ok, "expected a status error returned")
		require.Equal(t, http.StatusInternalServerError, serr.Code)
		require.Equal(t, "compliance overload", serr.Message)
	})

	t.Run("StatusError", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallInquiry, &openvasp.StatusError{Code: http.StatusBadRequest, Message: "high volume trading"})

		rep, err := client.Inquiry(inquiry)
		require.Nil(t, rep, "response should be nil in the case of an error")
		require.Error(t, err, "expected an error to be returned")

		serr, ok := err.(*openvasp.StatusError)
		require.True(t, ok, "expected a status error returned")
		require.Equal(t, http.StatusBadRequest, serr.Code)
		require.Equal(t, "high volume trading", serr.Message)
	})

	t.Run("StandardResponse", func(t *testing.T) {
		defer mock.Reset()
		mock.CallInquiry = func(*openvasp.Inquiry) (*openvasp.InquiryResolution, error) { return nil, nil }

		rep, err := client.Inquiry(inquiry)
		require.NoError(t, err, "did not expect an error")
		require.Equal(t, http.StatusOK, rep.StatusCode)

		info := rep.Info()
		require.Equal(t, inquiry.TRP, info, "expected incoming info to match outgoing info")

		resolution, err := rep.InquiryResolution()
		require.NoError(t, err, "could not parse inquiry resolution")
		require.Equal(t, openvasp.APIVersion, resolution.Version)
	})
}

func TestClientConfirmation(t *testing.T) {
	// Setup the mock handler the http test server.
	mock := &MockHandler{}
	handler := openvasp.TransferConfirmation(mock)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// Setup the client and the confirmation fixture.
	client := openvasp.NewClient()
	confirm := &openvasp.Confirmation{
		TRP: &openvasp.TRPInfo{
			Address:           srv.URL + "/x/12345?t=i",
			RequestIdentifier: "ba3bc19e-50f8-4e48-8a2a-da1e50f9b672",
		},
		TXID: "foo",
	}

	t.Run("Error", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallConfirmation, errors.New("chaos ensues"))

		rep, err := client.Confirmation(confirm)
		require.Nil(t, rep, "expected nil response on error")
		require.Error(t, err, "expected an error to be returned")

		serr, ok := err.(*openvasp.StatusError)
		require.True(t, ok, "expected a status error returned")
		require.Equal(t, http.StatusInternalServerError, serr.Code)
		require.Equal(t, "chaos ensues", serr.Message)
	})

	t.Run("Standard", func(t *testing.T) {
		defer mock.Reset()
		mock.UseError(CallConfirmation, nil)

		rep, err := client.Confirmation(confirm)
		require.NoError(t, err, "expected no error returned")
		require.Equal(t, http.StatusNoContent, rep.StatusCode)

		info := rep.Info()
		require.Equal(t, confirm.TRP, info)
	})

}
