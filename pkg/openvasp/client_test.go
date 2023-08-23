package openvasp_test

import (
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
