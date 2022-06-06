package keys_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/keys"
	"github.com/trisacrypto/trisa/pkg/trust"
)

func TestCertificateToExchange(t *testing.T) {
	// Tests the workflow where a Certificate is loaded, then serialized via Proto and
	// exchanged back to the user and loaded into an Exchange. This is the common TRISA
	// key exchange workflow and must work in all cases.

	// Load Certificate fixture with private keys
	sz, err := trust.NewSerializer(false)
	require.NoError(t, err, "could not create serializer to load fixture")

	provider, err := sz.ReadFile("testdata/certs.pem")
	require.NoError(t, err, "could not read test fixture")

	certs, err := keys.FromProvider(provider)
	require.NoError(t, err, "could not create Key from provider")
	require.True(t, certs.IsPrivate(), "expected test certs fixture to be private")

	// "exchange" keys with counterparty
	msg, err := certs.Proto()
	require.NoError(t, err, "could not create api.SigningKey message")

	xchange, err := keys.FromSigningKey(msg)
	require.NoError(t, err, "could not create exchange key from SigningKey message")
	require.False(t, xchange.IsPrivate(), "expected exchange keys not to be private")

	// Compare the certs with the xchange keys
	csk, err := certs.SealingKey()
	require.NoError(t, err, "could not extract sealing key from certs")
	xsk, err := xchange.SealingKey()
	require.NoError(t, err, "could not extract sealing key from exchange")
	require.Equal(t, csk, xsk, "expected sealing key to be identical")

	require.Equal(t, certs.PublicKeyAlgorithm(), xchange.PublicKeyAlgorithm(), "expected public key algorithms to match")

	cpks, err := certs.PublicKeySignature()
	require.NoError(t, err, "could not extract public key signature from certs")
	xpks, err := xchange.PublicKeySignature()
	require.NoError(t, err, "could not extract public key signature from exchange")
	require.Equal(t, cpks, xpks, "expected public key signatures to be identical")
}
