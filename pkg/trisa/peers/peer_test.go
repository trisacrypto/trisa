package peers_test

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that the ExchangeKeys function correctly retrieves a key from the endpoint.
func TestExchangeKeys(t *testing.T) {
	cache := makePeersCache(t)

	// Retrieve test peer
	p, err := cache.Get("test-peer")
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, "test-peer", p.Info().CommonName)
	require.NotNil(t, p.Info().SigningKey)

	// If signing key already exists, it should not be overwritten
	key, err := p.ExchangeKeys(false)
	require.NoError(t, err)
	require.Equal(t, &rsa.PublicKey{}, key)
	require.Equal(t, &rsa.PublicKey{}, p.Info().SigningKey)

	// Test concurrent ExchangeKeys calls
	t.Run("exchangeKeys", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
				t.Parallel()
				key, err := p.ExchangeKeys(true)
				require.NoError(t, err)
				require.NotNil(t, key)
				require.NotEqual(t, &rsa.PublicKey{}, p.Info().SigningKey)
			})
		}
	})

	// Signing key should be overwritten
	key = p.Info().SigningKey
	require.NotEqual(t, &rsa.PublicKey{}, key)
	data, err := x509.MarshalPKIXPublicKey(key)
	require.NoError(t, err)
	require.NotEmpty(t, data)
}
