package peers_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1/mock"
	"github.com/trisacrypto/trisa/pkg/trisa/keys"
	"github.com/trisacrypto/trisa/pkg/trisa/peers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Test that the ExchangeKeys function correctly retrieves a key from the endpoint.
func TestExchangeKeys(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Create a mock remote counter party to perform the key exchange with
	remote := mock.New(nil)
	defer remote.Shutdown()

	// Create a signing key for the local peer client
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Add the peer to the peers cache
	cache.Add(&peers.PeerInfo{
		ID:                  "36463211-d44b-40bd-88ae-6e8594c2b62c",
		RegisteredDirectory: "testdirectory.org",
		CommonName:          "test-peer",
		Endpoint:            "test-peer:4444",
		SigningKey:          &signingKey.PublicKey,
	})

	// Retrieve a peer without a directory lookup and connect it to the mock peer
	p, err := cache.Get("test-peer")
	require.NoError(t, err)
	require.Equal(t, "test-peer", p.Info().CommonName)
	require.NotNil(t, p.Info().SigningKey)

	p.Connect(
		grpc.WithContextDialer(remote.Channel().Dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// If signing key already exists, it should not be overwritten
	key, err := p.ExchangeKeys(false)
	require.NoError(t, err)
	require.Equal(t, p.Info().SigningKey, key)

	// TODO: test case where exchange keys returns an error
	// TODO: review other cases we should test with code coverage

	// Handle case where a key exchange is successful
	remote.OnKeyExchange = func(context.Context, *api.SigningKey) (*api.SigningKey, error) {
		certs, _, err := loadCertificates("testdata/server.pem")
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		rkey, err := keys.FromProvider(certs)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return rkey.Proto()
	}

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
	// TODO: please validate these checks are correctly checking what we expect
	key = p.Info().SigningKey
	require.NotEqual(t, signingKey.PublicKey, key)
	data, err := x509.MarshalPKIXPublicKey(key)
	require.NoError(t, err)
	require.NotEmpty(t, data)
}
