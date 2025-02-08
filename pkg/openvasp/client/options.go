package client

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/trisacrypto/trisa/pkg/trust"
)

// ClientOption allows us to configure the APIv1 client when it is created.
type ClientOption func(c *Client) error

// Specify the underlying HTTP client to use to make requests for complete control of
// how the HTTP transport is configured. This will overwrite any default configuration
// that the default client uses.
func WithClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

// Specify the mTLS or TLS configuration to use to connect to secure endpoints.
func WithTLSConfig(tlsConfig *tls.Config) ClientOption {
	return func(c *Client) error {
		c.client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		return nil
	}
}

// WithMTLS allows the user to specify an mTLS provider and pool for the client.
func WithMTLS(client *trust.Provider, pool trust.ProviderPool) ClientOption {
	return func(c *Client) (err error) {
		var crt tls.Certificate
		if crt, err = client.GetKeyPair(); err != nil {
			return err
		}

		var cas *x509.CertPool
		if cas, err = pool.GetCertPool(false); err != nil {
			return err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{crt},
			MinVersion:   tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
			RootCAs: cas,
		}

		c.client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		return nil
	}
}

// Specify the default API version to use in requests.
func WithAPIVersion(version string) ClientOption {
	return func(c *Client) error {
		c.apiVersion = version
		return nil
	}
}

// Specify the default API extensions to use in requests (overwrites default extensions).
func WithAPIExtensions(extensions ...string) ClientOption {
	return func(c *Client) error {
		c.extensions = extensions
		return nil
	}
}
