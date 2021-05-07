package trust

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
)

// PEMEncodePrivateKey as a PKCS8 ASN.1 DER key and write a PEM block with type "PRIVATE KEY"
func PEMEncodePrivateKey(key interface{}) ([]byte, error) {
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Type: BlockPrivateKey, Bytes: pkcs8}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// PEMDecodePrivateKey from a PEM encoded block. If the block type is "EC PRIVATE KEY",
// then the block is parsed as an EC private key in SEC 1, ASN.1 DER form. If the block
// is "RSA PRIVATE KEY" then it is decoded as a PKCS 1, ASN.1 DER form. If the block
// type is "PRIVATE KEY", the block is decoded as a PKCS 8 ASN.1 DER key, if that fails,
// then the PKCS 1 and EC parsers are tried in that order, before returning an error.
func PEMDecodePrivateKey(in []byte) (key interface{}, err error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, ErrDecodePrivateKey
	}

	return parsePrivateKey(block)
}

func parsePrivateKey(block *pem.Block) (key interface{}, err error) {
	// EC PRIVATE KEY specific handling
	if block.Type == BlockECPrivateKey {
		return x509.ParseECPrivateKey(block.Bytes)
	}

	// RSA PRIVATE KEY specific handling
	if block.Type == BlockRSAPrivateKey {
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	// Expect PRIVATE KEY if not EC or RSA at this point
	if block.Type != BlockPrivateKey {
		return nil, ErrDecodePrivateKey
	}

	// Try parsing private key using PKCS8, PKCS1, then EC
	if key, err = x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	if key, err = x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// Could not parse the private key
	return nil, ErrDecodePrivateKey
}

// PEMEncodeCertificate and write a PEM block with type "CERTIFICATE"
func PEMEncodeCertificate(c *x509.Certificate) ([]byte, error) {
	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Type: BlockCertificate, Bytes: c.Raw}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// PEMDecodeCertificate from PEM encoded block with type "CERTIFICATE"
func PEMDecodeCertificate(in []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(in)
	if block == nil || block.Type != BlockCertificate {
		return nil, ErrDecodeCertificate
	}
	return x509.ParseCertificate(block.Bytes)
}

// PEMEncodeCSR and write a PEM block with type "CERTIFICATE REQUEST"
func PEMEncodeCSR(c *x509.CertificateRequest) ([]byte, error) {
	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Type: BlockCertificateRequest, Bytes: c.Raw}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// PEMDecodeCSR from PEM encoded block with type "CERTIFICATE REQUEST"
func PEMDecodeCSR(in []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(in)
	if block == nil || block.Type != BlockCertificateRequest {
		return nil, ErrDecodeCSR
	}
	return x509.ParseCertificateRequest(block.Bytes)
}
