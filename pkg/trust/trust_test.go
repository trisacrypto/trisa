package trust_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trust"
	"github.com/trisacrypto/trisa/pkg/trust/mock"
	"software.sslmate.com/src/go-pkcs12"
)

func TestPrivateProvider(t *testing.T) {
	pfxData, err := mock.Chain()
	require.NoError(t, err)

	p, err := trust.Decrypt(pfxData, pkcs12.DefaultPassword)
	require.NoError(t, err)
	require.True(t, p.IsPrivate())
	require.Equal(t, "Test", p.String())

	pool, err := p.GetCertPool()
	require.NoError(t, err)
	require.Len(t, pool.Subjects(), 3)

	pair, err := p.GetKeyPair()
	require.NoError(t, err)
	require.NotNil(t, pair.PrivateKey)
	require.Len(t, pair.Certificate, 3)

	// Test encrypt/decrypt
	pfxData, err = p.Encrypt("supersecretsquirrel")
	require.NoError(t, err)

	_, err = trust.Decrypt(pfxData, "knockknock")
	require.Error(t, err)

	o, err := trust.Decrypt(pfxData, "supersecretsquirrel")
	require.NoError(t, err)
	require.Equal(t, p, o)
	require.True(t, o.IsPrivate())

	// Test encode/decode
	pfxData, err = p.Encode()
	require.NoError(t, err)

	o = &trust.Provider{}
	require.NotEqual(t, p, o)
	require.NoError(t, o.Decode(pfxData))
	require.Equal(t, p, o)
	require.True(t, o.IsPrivate())
}

func TestPublicProvider(t *testing.T) {
	pfxData, err := mock.Chain()
	require.NoError(t, err)

	priv, err := trust.Decrypt(pfxData, pkcs12.DefaultPassword)
	require.NoError(t, err)

	p := priv.Public()
	require.NotEqual(t, p, priv)

	o := p.Public()
	require.Equal(t, &p, &o)

	require.False(t, p.IsPrivate())
	require.Equal(t, "Test", p.String())

	pool, err := p.GetCertPool()
	require.NoError(t, err)
	require.Len(t, pool.Subjects(), 3)

	provPool := trust.NewPool(o)
	require.Equal(t, provPool[o.String()], o)
	require.False(t, o.IsPrivate())

	_, err = p.GetKeyPair()
	require.Error(t, err)

	// Test encrypt
	_, err = p.Encrypt("supersecretsquirrel")
	require.Error(t, err)

	// Test encode/decode
	pfxData, err = p.Encode()
	require.NoError(t, err)

	o = &trust.Provider{}
	require.NotEqual(t, p, o)
	require.NoError(t, o.Decode(pfxData))
	require.Equal(t, p, o)
	require.False(t, o.IsPrivate())
}
