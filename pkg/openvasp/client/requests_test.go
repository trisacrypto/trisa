package client_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp"
	"github.com/trisacrypto/trisa/pkg/openvasp/client"
	"github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

func TestNewRequest(t *testing.T) {
	ctx := context.Background()
	ta := &trp.Info{
		APIVersion: "v1",
		Address:    "https://example.com",
	}

	req, err := client.NewRequest(ctx, http.MethodGet, ta, nil)
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, ta.Address, req.URL.String())
	require.Equal(t, client.UserAgent, req.Header.Get("User-Agent"))
	require.Equal(t, client.Accept, req.Header.Get("Accept"))
	require.Equal(t, client.AcceptLanguage, req.Header.Get("Accept-Language"))
	require.Equal(t, client.AcceptEncode, req.Header.Get("Accept-Encoding"))
	require.Equal(t, "v1", req.Header.Get(openvasp.APIVersionHeader))
	require.NotEmpty(t, req.Header.Get(openvasp.RequestIdentifierHeader))
}

func TestNewRequestWithBody(t *testing.T) {
	ctx := context.Background()
	ta := &trp.Info{
		APIVersion: "v1",
		Address:    "https://example.com",
	}

	body := bytes.NewBufferString(`{"key":"value"}`)
	req, err := client.NewRequest(ctx, http.MethodPost, ta, body)
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, ta.Address, req.URL.String())
	require.Equal(t, openvasp.ContentTypeValue, req.Header.Get(openvasp.ContentTypeHeader))
}

func TestNewJSONRequest(t *testing.T) {
	ctx := context.Background()
	ta := &trp.Info{
		APIVersion: "v1",
		Address:    "https://example.com",
	}

	data := map[string]string{"key": "value"}
	c, _ := client.New()
	req, err := c.NewJSONRequest(ctx, http.MethodPost, ta, data)
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, openvasp.ContentTypeValue, req.Header.Get(openvasp.ContentTypeHeader))
}

func TestNewTextRequest(t *testing.T) {
	ctx := context.Background()
	ta := &trp.Info{
		APIVersion: "v1",
		Address:    "https://example.com",
	}

	data := "test string"
	c, _ := client.New()
	req, err := c.NewTextRequest(ctx, http.MethodPost, ta, data)
	require.NoError(t, err)
	require.NotNil(t, req)
	require.Equal(t, client.Accept, req.Header.Get("Accept"))
}
