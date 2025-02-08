package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/client"
	"github.com/trisacrypto/trisa/pkg/openvasp/extensions/discoverability"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

func TestNewRequest(t *testing.T) {
	t.Run("WithoutBody", func(t *testing.T) {
		ctx := context.Background()
		ta := &trp.Info{
			APIVersion:    "v1",
			APIExtensions: []string{"ivms101-extended", "secure-trisa-envelope"},
			Address:       "https://example.com",
		}

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		require.Equal(t, ta.Address, req.URL.String())

		expected := map[string]string{
			openvasp.APIVersionHeader:    "v1",
			openvasp.APIExtensionsHeader: "ivms101-extended, secure-trisa-envelope",
		}
		CheckHeaders(t, req, expected)

		// With no body, the content type header should be empty
		require.Empty(t, req.Header.Get(openvasp.ContentTypeHeader))
	})

	t.Run("WithBody", func(t *testing.T) {
		ctx := context.Background()
		ta := &trp.Info{
			APIVersion:        "v1",
			Address:           "https://example.com",
			RequestIdentifier: "748d4a01-01b7-4a50-8385-ce9d37dded49",
		}

		body := bytes.NewBufferString(`{"key":"value"}`)
		req, err := client.NewRequest(ctx, http.MethodPost, ta, body)
		require.NoError(t, err)
		require.NotNil(t, req)

		require.Equal(t, ta.Address, req.URL.String())

		expected := map[string]string{
			openvasp.APIVersionHeader:        "v1",
			openvasp.ContentTypeHeader:       openvasp.ContentTypeValue,
			openvasp.RequestIdentifierHeader: "748d4a01-01b7-4a50-8385-ce9d37dded49",
		}
		CheckHeaders(t, req, expected)
	})

	t.Run("WithContext", func(t *testing.T) {
		ctx := context.Background()
		ctx = client.ContextWithRequestID(ctx, "748d4a01-01b7-4a50-8385-ce9d37dded49")
		ctx = client.ContextWithTracing(ctx, "059f9071-a1e8-42ff-a9e6-89fe4f033593")

		ta := &trp.Info{Address: "https://example.com"}
		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		expected := map[string]string{
			openvasp.RequestIdentifierHeader: "748d4a01-01b7-4a50-8385-ce9d37dded49",
			"X-Request-ID":                   "059f9071-a1e8-42ff-a9e6-89fe4f033593",
		}
		CheckHeaders(t, req, expected)
	})

	t.Run("TAOverridesContext", func(t *testing.T) {
		ctx := context.Background()
		ctx = client.ContextWithRequestID(ctx, "748d4a01-01b7-4a50-8385-ce9d37dded49")

		ta := &trp.Info{
			Address:           "https://example.com",
			RequestIdentifier: "059f9071-a1e8-42ff-a9e6-89fe4f033593",
		}

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err)

		expected := map[string]string{
			openvasp.RequestIdentifierHeader: "059f9071-a1e8-42ff-a9e6-89fe4f033593",
		}
		CheckHeaders(t, req, expected)
	})

	t.Run("Defaults", func(t *testing.T) {
		ctx := context.Background()
		ta := &trp.Info{
			Address: "https://example.com",
		}

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err)
		require.NotNil(t, req)

		require.Equal(t, ta.Address, req.URL.String())
		CheckHeaders(t, req, nil)
	})

	t.Run("BadURL", func(t *testing.T) {
		ctx := context.Background()
		ta := &trp.Info{
			Address: "tanotatraveladdress",
		}

		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.Error(t, err)
		require.Nil(t, req)
	})
}

func TestClientRequests(t *testing.T) {
	client, err := client.New(
		client.WithAPIExtensions("ivms101-extended", "secure-trisa-envelope"),
		client.WithAPIVersion("2.3.1"),
	)

	require.NoError(t, err, "could not create client")
	require.Equal(t, "2.3.1", client.APIVersion())
	require.Equal(t, []string{"ivms101-extended", "secure-trisa-envelope"}, client.APIExtensions())

	ctx := context.Background()
	ta := &trp.Info{
		Address: "https://example.com",
	}

	t.Run("ClientDefaults", func(t *testing.T) {
		req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
		require.NoError(t, err)
		require.Equal(t, ta.Address, req.URL.String())

		expected := map[string]string{
			openvasp.APIVersionHeader:    "2.3.1",
			openvasp.APIExtensionsHeader: "ivms101-extended, secure-trisa-envelope",
		}
		CheckHeaders(t, req, expected)
	})

	t.Run("JSON", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		req, err := client.NewJSONRequest(ctx, http.MethodPost, ta, data)
		require.NoError(t, err)
		require.Equal(t, ta.Address, req.URL.String())

		expected := map[string]string{
			openvasp.APIVersionHeader:    "2.3.1",
			openvasp.APIExtensionsHeader: "ivms101-extended, secure-trisa-envelope",
			openvasp.ContentTypeHeader:   openvasp.ContentTypeValue,
		}
		CheckHeaders(t, req, expected)

		cmpt := make(map[string]string)
		err = json.NewDecoder(req.Body).Decode(&cmpt)
		require.NoError(t, err)
		require.Equal(t, data, cmpt)
	})

	t.Run("Text", func(t *testing.T) {
		expected := map[string]string{
			openvasp.APIVersionHeader:    "2.3.1",
			openvasp.APIExtensionsHeader: "ivms101-extended, secure-trisa-envelope",
			openvasp.ContentTypeHeader:   openvasp.MIMEPlainText,
			"Accept":                     openvasp.MIMEPlainText,
		}

		t.Run("String", func(t *testing.T) {
			data := "test string"
			req, err := client.NewTextRequest(ctx, http.MethodPost, ta, data)
			require.NoError(t, err)
			require.Equal(t, ta.Address, req.URL.String())
			CheckHeaders(t, req, expected)

			out, _ := io.ReadAll(req.Body)
			require.Equal(t, data, string(out))
		})

		t.Run("TextMarshaler", func(t *testing.T) {
			data := discoverability.Uptime(20 * time.Second)
			req, err := client.NewTextRequest(ctx, http.MethodPost, ta, data)
			require.NoError(t, err)
			require.Equal(t, ta.Address, req.URL.String())
			CheckHeaders(t, req, expected)

			out, _ := io.ReadAll(req.Body)
			require.Equal(t, "20", string(out))
		})

		t.Run("Bytes", func(t *testing.T) {
			data := []byte("foo")
			req, err := client.NewTextRequest(ctx, http.MethodPost, ta, data)
			require.NoError(t, err)
			require.Equal(t, ta.Address, req.URL.String())
			CheckHeaders(t, req, expected)

			out, _ := io.ReadAll(req.Body)
			require.Equal(t, "foo", string(out))
		})
	})
}

func TestDefaultHeaders(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		headers := map[string]string{}
		DefaultHeaders(t, headers)

		require.Len(t, headers, len(defaultHeaders))
		for key, val := range defaultHeaders {
			require.Equal(t, val, headers[key])
		}
	})

	t.Run("Nil", func(t *testing.T) {
		headers := DefaultHeaders(t, nil)
		require.Len(t, headers, len(defaultHeaders))
		for key, val := range defaultHeaders {
			require.Equal(t, val, headers[key])
		}
	})
}

func CheckHeaders(tb testing.TB, req *http.Request, expected map[string]string) {
	tb.Helper()
	expected = DefaultHeaders(tb, expected)
	for key, value := range expected {
		require.NotEmpty(tb, req.Header.Get(key), "missing header %s", key)
		require.Equal(tb, value, req.Header.Get(key), "incorrect header %s: got %s want %s", key, req.Header.Get(key), value)
	}

	require.NotEmpty(tb, req.Header.Get(openvasp.RequestIdentifierHeader), "no request identifier header")
}

var defaultHeaders = map[string]string{
	"User-Agent":              client.UserAgent,
	"Accept":                  client.Accept,
	"Accept-Language":         client.AcceptLanguage,
	"Accept-Encoding":         client.AcceptEncode,
	openvasp.APIVersionHeader: openvasp.APIVersion,
}

func DefaultHeaders(tb testing.TB, headers map[string]string) map[string]string {
	tb.Helper()
	if headers == nil {
		headers = make(map[string]string, len(defaultHeaders))
	}

	for key, value := range defaultHeaders {
		if _, ok := headers[key]; !ok {
			headers[key] = value
		}
	}

	return headers
}
