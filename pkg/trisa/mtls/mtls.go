package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"

	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config returns the standard TLS configuration for the TRISA network, loading the
// certificate from the specified provider. Using this TLS configuration ensures that
// all TRISA peer-to-peer connections are handled and verified correctly.
func Config(server *trust.Provider, clients trust.ProviderPool) (_ *tls.Config, err error) {
	if !server.IsPrivate() {
		return nil, errors.New("server provider must contain a private key to initialize TLS certs")
	}

	var crt tls.Certificate
	if crt, err = server.GetKeyPair(); err != nil {
		return nil, err
	}

	var pool *x509.CertPool
	if pool, err = clients.GetCertPool(false); err != nil {
		return nil, err
	}

	return &tls.Config{
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
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  pool,
	}, nil
}

// ServerCreds returns the grpc.ServerOption to create a gRPC server with mTLS.
func ServerCreds(server *trust.Provider, clients trust.ProviderPool) (_ grpc.ServerOption, err error) {
	var conf *tls.Config
	if conf, err = Config(server, clients); err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(conf)
	return grpc.Creds(creds), nil
}

// ClientCreds returns the grpc.DialOption to create a gRPC client with mTLS.
func ClientCreds(endpoint string, client *trust.Provider, servers trust.ProviderPool) (_ grpc.DialOption, err error) {
	if !client.IsPrivate() {
		return nil, errors.New("client provider must contain a private key to initialize TLS certs")
	}

	var crt tls.Certificate
	if crt, err = client.GetKeyPair(); err != nil {
		return nil, err
	}

	var u *url.URL
	if u, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("invalid endpoint: %q", err)
	}

	var pool *x509.CertPool
	if pool, err = servers.GetCertPool(false); err != nil {
		return nil, err
	}

	conf := &tls.Config{
		ServerName:   u.Host,
		Certificates: []tls.Certificate{crt},
		RootCAs:      pool,
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(conf)), nil
}
