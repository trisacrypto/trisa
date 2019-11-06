package trust

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
)

type Provider struct {
	chain tls.Certificate
}

func NewProvider(chain []byte) *Provider {
	p := &Provider{}
	p.AddChain(chain)
	return p
}

func (p *Provider) AddChain(in []byte) {
	var block *pem.Block
	for {
		block, in = pem.Decode(in)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			p.chain.Certificate = append(p.chain.Certificate, block.Bytes)
		}
	}
}

func (p *Provider) GetCertPool() *x509.CertPool {
	pool := x509.NewCertPool()
	for _, c := range p.chain.Certificate {
		x509Cert, err := x509.ParseCertificate(c)
		if err != nil {
			panic(err)
		}
		pool.AddCert(x509Cert)
	}
	return pool
}
