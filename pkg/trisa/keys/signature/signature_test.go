package signature_test

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/keys/signature"
)

func Example() {
	// Generate a new RSA key pair
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Sign the public key
	pks, _ := signature.New(&key.PublicKey)

	// Make sure the pks matches the original public key
	match := signature.Match(pks, &key.PublicKey)
	fmt.Printf("%t\n", match)
	// Output:
	// true
}

func TestRSA(t *testing.T) {
	pub, err := loadFixture(fixtureRSA)
	require.NoError(t, err, "could not load RSA fixtures")

	// Test default workflow
	pks, err := signature.New(pub)
	require.NoError(t, err, "could not create default signature")
	require.Equal(t, "SHA256:BbUdmS0/gwFi7loGUcJeEUD4REaT5XgtCoRX+BDnkNg", pks)
	require.True(t, signature.Match(pks, pub))

	// Test Valid signatures
	testCases := []struct {
		algorithm signature.Algorithm
		expected  string
	}{
		{signature.MD5, "MD5:oSbNMVJNnN/l4KgGVU+C2w"},
		{signature.SHA256, "SHA256:BbUdmS0/gwFi7loGUcJeEUD4REaT5XgtCoRX+BDnkNg"},
		{signature.SHA512, "SHA512:5xSYB7LfujoDzyzWJYqWDk7ffMsU8xqLCKbDnBcuLO75VF5fczy47AR9APMu0HVQrhKeejxeHV7YXj/702MGAQ"},
	}

	for _, testCase := range testCases {
		pks, err := signature.Sign(pub, testCase.algorithm)
		require.NoError(t, err, "could not sign using %s", testCase.algorithm.String())
		require.Equal(t, testCase.expected, pks)
		require.True(t, signature.Match(pks, pub))
	}

	// Test invalid signatures
	require.False(t, signature.Match("FOO", pub))
	require.False(t, signature.Match("FOO:oSbNMVJNnN/l4KgGVU+C2w", pub))
	require.False(t, signature.Match("MD5:F", pub))
}

func TestECDSA(t *testing.T) {
	pub, err := loadFixture(fixtureECDSA)
	require.NoError(t, err, "could not load ECDSA fixtures")

	// Test default workflow
	pks, err := signature.New(pub)
	require.NoError(t, err, "could not create default signature")
	require.Equal(t, "SHA256:KoKS3i9txi8DPwgmIdDr+rPrdF0bh5SbaQi4/HYheRA", pks)
	require.True(t, signature.Match(pks, pub))

	// Test Valid signatures
	testCases := []struct {
		algorithm signature.Algorithm
		expected  string
	}{
		{signature.MD5, "MD5:ypCkp1+OdN4sVYyvVeu3uQ"},
		{signature.SHA256, "SHA256:KoKS3i9txi8DPwgmIdDr+rPrdF0bh5SbaQi4/HYheRA"},
		{signature.SHA512, "SHA512:X68VhfJAUVFRqpmz0N6Wftp+JNvTZvO1gAb8IuOXLe9eG+BCosUfvXgAPD8TrEH0L0rVBZUYuVAAezjAsqLZfw"},
	}

	for _, testCase := range testCases {
		pks, err := signature.Sign(pub, testCase.algorithm)
		require.NoError(t, err, "could not sign using %s", testCase.algorithm.String())
		require.Equal(t, testCase.expected, pks)
		require.True(t, signature.Match(pks, pub))
	}

	// Test invalid signatures
	require.False(t, signature.Match("FOO", pub))
	require.False(t, signature.Match("FOO:ypCkp1+OdN4sVYyvVeu3uQ", pub))
	require.False(t, signature.Match("MD5:F", pub))
}

func TestED25519(t *testing.T) {
	pub, err := loadFixture(fixtureED25519)
	require.NoError(t, err, "could not load ED25519 fixtures")

	// Test default workflow
	pks, err := signature.New(pub)
	require.NoError(t, err, "could not create default signature")
	require.Equal(t, "SHA256:csSenMcPrKDqVgRW4JOjfFyOuZta+bPBLqWI4/cJKyU", pks)
	require.True(t, signature.Match(pks, pub))

	// Test Valid signatures
	testCases := []struct {
		algorithm signature.Algorithm
		expected  string
	}{
		{signature.MD5, "MD5:ckldayqiExTtcPlZBT8HCQ"},
		{signature.SHA256, "SHA256:csSenMcPrKDqVgRW4JOjfFyOuZta+bPBLqWI4/cJKyU"},
		{signature.SHA512, "SHA512:+s9qX6DWjnO/851NI69ViNkhD3iDQn36JG3NlZTDB/BWcppZf0aZ4AaHsIe+dpRXjiDKJN719x0J052//YYu5g"},
	}

	for _, testCase := range testCases {
		pks, err := signature.Sign(pub, testCase.algorithm)
		require.NoError(t, err, "could not sign using %s", testCase.algorithm.String())
		require.Equal(t, testCase.expected, pks)
		require.True(t, signature.Match(pks, pub))
	}

	// Test invalid signatures
	require.False(t, signature.Match("FOO", pub))
	require.False(t, signature.Match("FOO:ckldayqiExTtcPlZBT8HCQ", pub))
	require.False(t, signature.Match("MD5:F", pub))
}

const (
	fixtureRSA     = "testdata/rsapub.pem"
	fixtureECDSA   = "testdata/ecdsa.pem"
	fixtureED25519 = "testdata/ed25519.pem"
)

// Load a fixture to use in testing
func loadFixture(path string) (pub interface{}, err error) {
	// Check if path exists, generate fixture if it doesn't
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = generateFixture(path); err != nil {
			return nil, err
		}
	}

	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("could not decode PEM data from file")
	}
	return x509.ParsePKIXPublicKey(block.Bytes)
}

// If the fixtures aren't in testdata folder, generate them.
func generateFixture(path string) (err error) {
	switch path {
	case fixtureRSA:
		return generateRSAFixture(path)
	case fixtureECDSA:
		return generateECDSAFixture(path)
	case fixtureED25519:
		return generateED25519Fixture(path)
	default:
		return fmt.Errorf("unknown fixture %q", path)
	}
}

func generateRSAFixture(path string) (err error) {
	var key *rsa.PrivateKey
	if key, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return fmt.Errorf("could not create RSA key fixture: %s", err)
	}

	if err = writeKeyFixture(&key.PublicKey, path, "RSA PUBLIC KEY"); err != nil {
		return fmt.Errorf("could not create RSA key fixture: %s", err)
	}

	return nil
}

func generateECDSAFixture(path string) (err error) {
	var key *ecdsa.PrivateKey
	if key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		return fmt.Errorf("could not create ECDSA key fixture: %s", err)
	}

	if err = writeKeyFixture(&key.PublicKey, path, "ECDSA PUBLIC KEY"); err != nil {
		return fmt.Errorf("could not create ECDSA key fixture: %s", err)
	}

	return nil
}

func generateED25519Fixture(path string) (err error) {
	var key ed25519.PublicKey
	if key, _, err = ed25519.GenerateKey(rand.Reader); err != nil {
		return fmt.Errorf("could not create ED25519 key fixture: %s", err)
	}

	if err = writeKeyFixture(key, path, "ED25519 PUBLIC KEY"); err != nil {
		return fmt.Errorf("could not create ED25519 key fixture: %s", err)
	}

	return nil
}

func writeKeyFixture(pub interface{}, path, block string) (err error) {
	var pkix []byte
	if pkix, err = x509.MarshalPKIXPublicKey(pub); err != nil {
		return err
	}

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}

	if err = pem.Encode(f, &pem.Block{Type: block, Bytes: pkix}); err != nil {
		return err
	}
	return nil
}
