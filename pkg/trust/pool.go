package trust

import (
	"crypto/x509"
	"fmt"
)

// ProviderPool is a collection of provider objects, used to collectively manage the
// public trust certificates of peers in the TRISA network. The Pool maps common names
// to the Provider certificates and ensures that only public providers without private
// keys are stored in the pool.
type ProviderPool map[string]*Provider

// NewPool from the specified providers (optional).
func NewPool(providers ...*Provider) ProviderPool {
	pool := make(ProviderPool)
	for _, p := range providers {
		pool.Add(p)
	}
	return pool
}

// Add a provider to the pool, mapping its common name to the provider.
func (pool ProviderPool) Add(p *Provider) {
	pool[p.String()] = p.Public()
}

// GetCertPool returns the x509.CertPool certificates of all provider certificates with
// the intermediate and root CA certificates. If requested, the pool can also load the
// system certificates pool if any system certificates exit.
func (pool ProviderPool) GetCertPool(includeSystem bool) (_ *x509.CertPool, err error) {
	var certPool *x509.CertPool
	if includeSystem {
		if certPool, _ = x509.SystemCertPool(); certPool == nil {
			certPool = x509.NewCertPool()
		}
	} else {
		certPool = x509.NewCertPool()
	}

	for name, provider := range pool {
		for i, asn1Data := range provider.chain.Certificate {
			var cert *x509.Certificate
			if cert, err = x509.ParseCertificate(asn1Data); err != nil {
				return nil, fmt.Errorf("could not parse certificate %d of %s: %s", i, name, err)
			}
			certPool.AddCert(cert)
		}
	}

	return certPool, nil
}
