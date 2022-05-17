package keys

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/trisacrypto/trisa/pkg/trust"
)

// ParseKeyExchangeData attempts to parse key exchange data that may arrive in multiple
// seralized formats. In order, it first attempts to parse marshaled PKIX public keys,
// then PEM encoded data (either x509 certificates or public keys), finally attempting
// to parse raw x509 data. If more exchange formats are detected, they should be added
// to this methods parsing interface.
func ParseKeyExchangeData(data []byte) (pubkey interface{}, err error) {
	if pubkey, err = x509.ParsePKIXPublicKey(data); err == nil {
		return pubkey, nil
	}

	if pubkey, err = parsePEMData(data); err == nil {
		return pubkey, nil
	}

	// TODO: add trust.Serializer here to handle GDS data

	var cert *x509.Certificate
	if cert, err = x509.ParseCertificate(data); err == nil {
		return cert.PublicKey, nil
	}

	return nil, ErrUnparsableKeyExchange
}

func parsePEMData(data []byte) (_ interface{}, err error) {
	chain := make([]*x509.Certificate, 0)
	keys := make([]interface{}, 0)

	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil {
			break
		}

		switch block.Type {
		case trust.BlockCertificate:
			var cert *x509.Certificate
			if cert, err = x509.ParseCertificate(block.Bytes); err != nil {
				return nil, fmt.Errorf("could not parse certificate: %s", err)
			}
			chain = append(chain, cert)
		case trust.BlockPublicKey:
			var pubkey interface{}
			if pubkey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
				return nil, fmt.Errorf("could not parse pkix public key: %s", err)
			}
			keys = append(keys, pubkey)
		case trust.BlockRSAPublicKey:
			var pubkey interface{}
			if pubkey, err = x509.ParsePKCS1PublicKey(block.Bytes); err != nil {
				if pubkey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
					return nil, fmt.Errorf("could not parse rsa public key: %s", err)
				}
			}
			keys = append(keys, pubkey)
		}

	}

	// If we have both certs and keys we cannot identify which key to use.
	if len(chain) > 0 && len(keys) > 0 {
		return nil, ErrTooManyBlocks
	}

	// If we have keys but no certs, we should have at most one key
	if len(chain) == 0 && len(keys) > 0 {
		if len(keys) > 1 {
			return nil, ErrTooManyBlocks
		}
		return keys[0], nil
	}

	// If we have certs but no keys, we should return the first cert, which we expect
	// is the leaf certificate in a trust chain (but ideally, it's just one cert).
	if len(chain) > 0 && len(keys) == 0 {
		return chain[0].PublicKey, nil
	}

	// We have neither certs nor keys
	return nil, ErrNoPublicKey
}
