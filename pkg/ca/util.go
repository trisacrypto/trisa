package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

func PEMEncodePrivateKey(key interface{}) ([]byte, error) {
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func PEMEncodeCertificate(c *x509.Certificate) ([]byte, error) {
	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func PEMDecodeCertificate(in []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, fmt.Errorf("error decoding pem certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func PEMEncodeCSR(c *x509.CertificateRequest) ([]byte, error) {
	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: c.Raw}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func PEMDecodeCSR(in []byte) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode(in)
	if block == nil {
		return nil, fmt.Errorf("error decoding pem certificate")
	}
	return x509.ParseCertificateRequest(block.Bytes)
}

func GenerateRSAPrivateKey(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func CreateCSR(subject pkix.Name, key interface{}) ([]byte, error) {
	tpl := &x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.SHA512WithRSA,
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, tpl, key)
	if err != nil {
		return nil, err
	}
	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		return nil, err
	}
	return PEMEncodeCSR(csr)
}

func Sha256Fingerprint(c *x509.Certificate) string {
	sum := sha256.Sum256(c.Raw)
	return hex.EncodeToString(sum[:])
}
