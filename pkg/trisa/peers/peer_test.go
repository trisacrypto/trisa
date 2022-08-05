package peers_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
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

	// Handle the case where the key exchange returns an error
	require.NoError(t, remote.UseError(mock.KeyExchangeRPC, codes.Internal, "key exchange error"))

	// If signing key already exists, it should not be overwritten
	key, err := p.ExchangeKeys(false)
	require.NoError(t, err)
	require.Equal(t, p.Info().SigningKey, key)

	// Should return an error if key exchange is forced but the RPC returns an error
	_, err = p.ExchangeKeys(true)
	require.Error(t, err, "expected exchange keys to return an error")
	require.Nil(t, p.SigningKey(), "expected signing key to be nil")

	// Should return an error if the key exchange returns an unparseable key
	remote.OnKeyExchange = func(context.Context, *api.SigningKey) (*api.SigningKey, error) {
		return &api.SigningKey{
			Data: []byte("not a key"),
		}, nil
	}
	_, err = p.ExchangeKeys(true)
	require.Error(t, err, "expected exchange keys with invalid key to return an error")

	// Should return an error if the key exchange returns a key in the wrong format
	remote.OnKeyExchange = func(context.Context, *api.SigningKey) (*api.SigningKey, error) {
		var (
			key  *ecdsa.PrivateKey
			data []byte
			err  error
		)
		if key, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader); err != nil {
			return nil, err
		}
		if data, err = x509.MarshalPKIXPublicKey(&key.PublicKey); err != nil {
			return nil, err
		}
		return &api.SigningKey{
			Data: data,
		}, nil
	}
	_, err = p.ExchangeKeys(true)
	require.EqualError(t, err, fmt.Sprintf("unsupported public key type %T", &ecdsa.PublicKey{}))

	// Load the certificate fixtures
	certs, _, err := loadCertificates("testdata/server.pem")
	require.NoError(t, err, "could not load certificate fixtures")
	privateKey, err := certs.GetRSAKeys()
	require.NoError(t, err, "could not extract rsa keys from fixtures")
	publicKey := &privateKey.PublicKey
	publicData, err := x509.MarshalPKIXPublicKey(publicKey)
	require.NoError(t, err, "could not marshal public key from fixtures")

	// Handle case where a key exchange is successful
	remote.OnKeyExchange = func(context.Context, *api.SigningKey) (*api.SigningKey, error) {
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
	key = p.Info().SigningKey
	require.Equal(t, publicKey, key)
	data, err := x509.MarshalPKIXPublicKey(key)
	require.NoError(t, err)
	require.Equal(t, publicData, data)
}
