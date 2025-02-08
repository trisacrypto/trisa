package client_test

import (
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/client"
	"github.com/trisacrypto/trisa/pkg/openvasp/extensions/discoverability"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

func TestClient(t *testing.T) {
	client, err := client.New()
	require.NoError(t, err, "could not create client")
	require.NotNil(t, client, "client should not be nil")

	ctx := context.Background()

	t.Run("Identity", func(t *testing.T) {
		srv, _ := NewServer(ValidatePath("/identity", HandleFixture(http.MethodGet, "testdata/identity.json")))
		defer srv.Close()

		id, err := client.Identity(ctx, srv.URL)
		require.NoError(t, err, "could not execute request")
		require.NotNil(t, id, "identity should not be nil")
		require.Equal(t, "Acme VASP", id.Name)
		require.Equal(t, "TVYD005I7Q5IJK0EMI53", id.LEI)
		require.True(t, strings.HasPrefix(id.X509, "-----BEGIN CERTIFICATE-----"))
	})

	t.Run("Inquiry", func(t *testing.T) {
		srv, ta := NewServer(HandleFixture(http.MethodPost, "testdata/acknowledged.json"))
		defer srv.Close()

		inquiry := &trp.Inquiry{}
		Fixture(t, "testdata/inquiry.json", inquiry)
		inquiry.Info = ta

		out, err := client.Inquiry(ctx, inquiry)
		require.NoError(t, err, "could not execute request")
		require.NotNil(t, out, "no resolution returned")
		require.Equal(t, "3.2.1", out.Version)
		require.Empty(t, out.Approved)
		require.Empty(t, out.Rejected)
	})

	t.Run("Resolve", func(t *testing.T) {
		srv, ta := NewServer(HandleNoContent(http.MethodPost))
		defer srv.Close()

		resolution := &trp.Resolution{}
		Fixture(t, "testdata/approved.json", resolution)
		resolution.Info = ta

		err := client.Resolve(ctx, resolution)
		require.NoError(t, err, "could not execute request")
	})

	t.Run("Confirm", func(t *testing.T) {
		srv, ta := NewServer(HandleNoContent(http.MethodPost))
		defer srv.Close()

		confirmation := &trp.Confirmation{}
		Fixture(t, "testdata/confirmation.json", confirmation)
		confirmation.Info = ta

		err := client.Confirm(ctx, confirmation)
		require.NoError(t, err, "could not execute request")
	})

	t.Run("Version", func(t *testing.T) {
		srv, _ := NewServer(ValidatePath("/version", HandleFixture(http.MethodGet, "testdata/version.json")))
		defer srv.Close()

		vers, err := client.Version(ctx, srv.URL)
		require.NoError(t, err, "could not execute request")
		require.NotNil(t, vers, "version should not be nil")
		require.Equal(t, "1.2.0", vers.Version)
		require.Equal(t, "21Analytics", vers.Vendor)
	})

	t.Run("Uptime", func(t *testing.T) {
		srv, _ := NewServer(ValidatePath("/uptime", HandleUptime(420*time.Second)))
		defer srv.Close()

		uptime, err := client.Uptime(ctx, srv.URL)
		require.NoError(t, err, "could not execute request")
		require.Equal(t, 420*time.Second, time.Duration(uptime))
	})

	t.Run("Extensions", func(t *testing.T) {
		srv, _ := NewServer(ValidatePath("/extensions", HandleFixture(http.MethodGet, "testdata/extensions.json")))
		defer srv.Close()

		extensions, err := client.Extensions(ctx, srv.URL)
		require.NoError(t, err, "could not execute request")
		require.NotNil(t, extensions, "extensions should not be nil")
		require.Equal(t, []string{"extended-ivms101", "deterministic-transfer"}, extensions.Supported)
		require.Equal(t, []string{"message-signing"}, extensions.Required)
	})

	t.Run("Error", func(t *testing.T) {
		srv, ta := NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "something bad happened", http.StatusUnprocessableEntity)
		}))
		defer srv.Close()

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err, "could not create request")

		rep, err := client.Do(req, nil)
		require.NotNil(t, rep, "expected reply to be returned")

		serr, ok := err.(*trp.StatusError)
		require.True(t, ok, "expected status error")
		require.Equal(t, http.StatusUnprocessableEntity, serr.Code)
		require.Equal(t, "something bad happened", serr.Message)
	})

	t.Run("ErrorCode", func(t *testing.T) {
		srv, ta := NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}))
		defer srv.Close()

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err, "could not create request")

		rep, err := client.Do(req, nil)
		require.NotNil(t, rep, "expected reply to be returned")

		serr, ok := err.(*trp.StatusError)
		require.True(t, ok, "expected status error")
		require.Equal(t, http.StatusTeapot, serr.Code)
		require.Equal(t, http.StatusText(http.StatusTeapot), serr.Message)
	})

	t.Run("ErrorJSON", func(t *testing.T) {
		srv, ta := NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", openvasp.MIMEJSON)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": 418, "message": "bad teapot"}`))
		}))
		defer srv.Close()

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err, "could not create request")

		rep, err := client.Do(req, nil)
		require.NotNil(t, rep, "expected reply to be returned")

		serr, ok := err.(*trp.StatusError)
		require.True(t, ok, "expected status error")
		require.Equal(t, http.StatusTeapot, serr.Code)
		require.Equal(t, "bad teapot", serr.Message)
	})

	t.Run("Do", func(t *testing.T) {
		t.Run("NoContent", func(t *testing.T) {
			srv, ta := NewServer(HandleNoContent(http.MethodGet))
			defer srv.Close()

			req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
			require.NoError(t, err, "could not create request")

			rep, err := client.Do(req, nil)
			require.NoError(t, err, "could not execute request")
			require.NotNil(t, rep, "expected reply to be returned")

			require.Nil(t, rep.Err(), "expected no error")
			serr, ok := rep.StatusError()
			require.Nil(t, serr, "expected no status error")
			require.True(t, ok, "expected status error to be nil")

			info := rep.Info()
			require.NotNil(t, info, "expected TRP info")
			require.Equal(t, ta.Address, info.Address)
			require.Equal(t, openvasp.APIVersion, info.APIVersion)
			require.Equal(t, ta.RequestIdentifier, info.RequestIdentifier)
			require.Empty(t, info.APIExtensions)
		})

		t.Run("Error", func(t *testing.T) {
			srv, ta := NewServer(ValidateAPIHeaders(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set(openvasp.APIExtensionsHeader, "extended-ivms101, deterministic-transfer")
				http.Error(w, "something bad happened", http.StatusUnprocessableEntity)
			}))
			defer srv.Close()

			ta.APIExtensions = []string{"extended-ivms101", "deterministic-transfer"}
			req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
			require.NoError(t, err, "could not create request")

			rep, err := client.Do(req, nil)
			require.Error(t, err, "expected a status error")
			require.NotNil(t, rep, "expected reply to be returned")

			serr, ok := rep.StatusError()
			require.True(t, ok, "expected status error")
			require.Equal(t, http.StatusUnprocessableEntity, serr.Code)
			require.Equal(t, "something bad happened", serr.Message)

			info := rep.Info()
			require.NotNil(t, info, "expected TRP info")
			require.Equal(t, ta.Address, info.Address)
			require.Equal(t, openvasp.APIVersion, info.APIVersion)
			require.Equal(t, ta.RequestIdentifier, info.RequestIdentifier)
			require.Equal(t, ta.APIExtensions, info.APIExtensions)
		})
	})
}

func NewServer(handler http.Handler) (*httptest.Server, *trp.Info) {
	srv := httptest.NewServer(handler)
	ta := &trp.Info{Address: srv.URL}
	return srv, ta
}

func HandleFixture(method, path string) http.HandlerFunc {
	return ValidateAPIHeaders(method, func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)
	})
}

func HandleUptime(dur time.Duration) http.HandlerFunc {
	return ValidateAPIHeaders(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		data, err := discoverability.Uptime(dur).MarshalText()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

func HandleNoContent(method string) http.HandlerFunc {
	return ValidateAPIHeaders(method, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func ValidatePath(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func ValidateAPIHeaders(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var apiVersion string
		if apiVersion = r.Header.Get(openvasp.APIVersionHeader); apiVersion == "" {
			http.Error(w, "missing API version header", http.StatusBadRequest)
			return
		}

		var requestIdentifier string
		if requestIdentifier = r.Header.Get(openvasp.RequestIdentifierHeader); requestIdentifier == "" {
			http.Error(w, "must specify request identifier", http.StatusBadRequest)
			return
		}

		// Set the APIVersion header in the outgoing response
		w.Header().Add(openvasp.APIVersionHeader, openvasp.APIVersion)

		// Echo back the request identifier in the outgoing response
		w.Header().Add(openvasp.RequestIdentifierHeader, requestIdentifier)

		// Enforce content type header
		if contentType := r.Header.Get(openvasp.ContentTypeHeader); contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "malformed content-type header", http.StatusUnsupportedMediaType)
				return
			}

			if mt != openvasp.MIMEJSON && mt != openvasp.MIMEPlainText {
				http.Error(w, "content-type header must be application/json or text/plain", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}

func Fixture(t testing.TB, path string, v interface{}) {
	t.Helper()

	f, err := os.Open(path)
	require.NoError(t, err, "could not open fixture file")
	defer f.Close()

	err = json.NewDecoder(f).Decode(v)
	require.NoError(t, err, "could not unmarshal fixture")
}
