package trust

import (
	"context"
	"crypto/x509"
	"time"

	"github.com/trisacrypto/trisa/pkg/trisa/discovery"
)

type Provider struct {
	caDisco   []*discovery.Trisa
	rootCAs   []*x509.Certificate
	issuerCAs []*x509.Certificate
}

// This needs improvement for refresh/reloads. Doing full init for now.
func NewProvider(caURLs []string) (*Provider, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := &Provider{}

	for _, url := range caURLs {

		// Load disco client per trusted CA
		disco, err := discovery.New(url)
		if err != nil {
			return nil, err
		}
		if err := disco.Init(ctx); err != nil {
			return nil, err
		}
		if err := disco.LoadAll(ctx); err != nil {
			return nil, err
		}

		// Store disco client for later refresh ...
		p.caDisco = append(p.caDisco, disco)

		// Load it all for ourselves for now
		p.rootCAs = append(p.rootCAs, disco.RootCAs...)
		p.issuerCAs = append(p.issuerCAs, disco.IssuerCAs...)
	}

	return p, nil
}

func (p *Provider) GetCertPool() *x509.CertPool {
	pool := x509.NewCertPool()
	for _, crt := range p.rootCAs {
		pool.AddCert(crt)
	}
	for _, crt := range p.issuerCAs {
		pool.AddCert(crt)
	}
	return pool
}
