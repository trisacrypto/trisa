package peers_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/require"
	gdsmock "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1/mock"
	"github.com/trisacrypto/trisa/pkg/trisa/peers"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

// Test that Add correctly adds peers to the Peers cache.
func TestAdd(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Common name is required
	require.Error(t, cache.Add(&peers.PeerInfo{}))

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Test adding peers concurrently
	t.Run("addTests", func(t *testing.T) {
		tests := []struct {
			name string
			info *peers.PeerInfo
		}{
			{"id", &peers.PeerInfo{
				CommonName: "leonardo",
				ID:         "1",
			}},
			{"directory", &peers.PeerInfo{
				CommonName:          "leonardo",
				RegisteredDirectory: "testdirectory.org",
			}},
			{"endpoint", &peers.PeerInfo{
				CommonName: "leonardo",
				Endpoint:   "https://leonardo.trisatest.net:443",
			}},
			{"key", &peers.PeerInfo{
				CommonName: "leonardo",
				SigningKey: &privateKey.PublicKey,
			}},
			{"differentPeer", &peers.PeerInfo{
				CommonName:          "donatello",
				ID:                  "2",
				RegisteredDirectory: "testdirectory.org",
				Endpoint:            "https://donatello.trisatest.net:443",
			}},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				require.NoError(t, cache.Add(tt.info))
			})
		}
	})

	// Verify the final state of the cache
	leonardo, err := cache.Get("leonardo")
	require.NoError(t, err)
	require.Equal(t, "leonardo", leonardo.Info().CommonName)
	require.Equal(t, "1", leonardo.Info().ID)
	require.Equal(t, "testdirectory.org", leonardo.Info().RegisteredDirectory)
	require.Equal(t, "https://leonardo.trisatest.net:443", leonardo.Info().Endpoint)
	require.Equal(t, &privateKey.PublicKey, leonardo.Info().SigningKey)

	donatello, err := cache.Get("donatello")
	require.NoError(t, err)
	require.Equal(t, "donatello", donatello.Info().CommonName)
	require.Equal(t, "2", donatello.Info().ID)
	require.Equal(t, "testdirectory.org", donatello.Info().RegisteredDirectory)
	require.Equal(t, "https://donatello.trisatest.net:443", donatello.Info().Endpoint)
}

// Test that FromContext returns the correct Peer given the connection context.
func TestFromContext(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Add a peer to the cache
	require.NoError(t, cache.Add(&peers.PeerInfo{
		CommonName: "leonardo",
		ID:         "1",
	}))

	// Context does not contain a peer
	ctx := context.Background()
	_, err = cache.FromContext(ctx)
	require.Error(t, err)

	// Peer has badly formatted credentials
	remotePeer := peer.Peer{
		AuthInfo: nil,
	}
	_, err = cache.FromContext(peer.NewContext(ctx, &remotePeer))
	require.Error(t, err)

	// Peer has no certificate
	remotePeer.AuthInfo = credentials.TLSInfo{}
	_, err = cache.FromContext(peer.NewContext(ctx, &remotePeer))
	require.Error(t, err)

	remotePeer.AuthInfo = credentials.TLSInfo{
		State: tls.ConnectionState{
			VerifiedChains: nil,
		},
	}
	_, err = cache.FromContext(peer.NewContext(ctx, &remotePeer))
	require.Error(t, err)

	remotePeer.AuthInfo = credentials.TLSInfo{
		State: tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{{}},
		},
	}
	_, err = cache.FromContext(peer.NewContext(ctx, &remotePeer))
	require.Error(t, err)

	// Certificate has no common name
	remotePeer.AuthInfo = credentials.TLSInfo{
		State: tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{
				{
					{
						Subject: pkix.Name{},
					},
				},
			},
		},
	}
	_, err = cache.FromContext(peer.NewContext(ctx, &remotePeer))
	require.Error(t, err)

	// Peer does not exist in the cache - should be created
	remotePeer.AuthInfo = credentials.TLSInfo{
		State: tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{
				{
					{
						Subject: pkix.Name{
							CommonName: "donatello",
						},
					},
				},
			},
		},
	}

	// Test calling FromContext concurrently
	t.Run("fromContext", func(t *testing.T) {
		tests := []struct {
			name string
			Peer string
		}{
			{"alreadyExists", "leonardo"},
			{"alreadyExists2", "leonardo"},
			{"newPeer", "donatello"},
			{"newPeer2", "donatello"},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				remotePeer := peer.Peer{
					AuthInfo: credentials.TLSInfo{
						State: tls.ConnectionState{
							VerifiedChains: [][]*x509.Certificate{
								{
									{
										Subject: pkix.Name{
											CommonName: tt.Peer,
										},
									},
								},
							},
						},
					},
				}
				p, err := cache.FromContext(peer.NewContext(ctx, &remotePeer))
				require.NoError(t, err)
				require.Equal(t, tt.Peer, p.Info().CommonName)
			})
		}
	})

	// Cache should contain the two peers
	leonardo, err := cache.Get("leonardo")
	require.NoError(t, err)
	require.Equal(t, "leonardo", leonardo.Info().CommonName)
	require.Equal(t, "1", leonardo.Info().ID)

	donatello, err := cache.Get("donatello")
	require.NoError(t, err)
	require.Equal(t, "donatello", donatello.Info().CommonName)
}

// Test that the Lookup function returns the correct remote peer given the common name.
func TestLookup(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Remote peer does not exist in the directory
	_, err = cache.Lookup("missing")
	require.Error(t, err)

	// Test concurrent Lookup calls
	t.Run("lookup", func(t *testing.T) {
		tests := []struct {
			name string
			peer string
		}{
			{"validCert", "leonardo"},
			{"validCert2", "leonardo"},
			{"invalidCert", "donatello"},
			{"invalidCert2", "donatello"},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				p, err := cache.Lookup(tt.peer)
				require.NoError(t, err)
				require.NotNil(t, p)
				// require.Equal(t, apimock.RemotePeers[tt.peer].ID, p.Info().ID)
			})
		}
	})

	// Cache should contain the two peers
	leonardo, err := cache.Get("leonardo")
	require.NoError(t, err)
	// expected := apimock.RemotePeers["leonardo"]
	// require.Equal(t, expected.ID, leonardo.Info().ID)
	// require.Equal(t, expected.RegisteredDirectory, leonardo.Info().RegisteredDirectory)
	// require.Equal(t, expected.CommonName, leonardo.Info().CommonName)
	// require.Equal(t, expected.Endpoint, leonardo.Info().Endpoint)
	require.NotNil(t, leonardo.Info().SigningKey)

	donatello, err := cache.Get("donatello")
	require.NoError(t, err)
	// expected = apimock.RemotePeers["donatello"]
	// require.Equal(t, expected.ID, donatello.Info().ID)
	// require.Equal(t, expected.RegisteredDirectory, donatello.Info().RegisteredDirectory)
	// require.Equal(t, expected.CommonName, donatello.Info().CommonName)
	// require.Equal(t, expected.Endpoint, donatello.Info().Endpoint)
	require.Nil(t, donatello.Info().SigningKey)
}

// Test that the Search function returns the matching remote peer given the name.
func TestSearch(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Remote peer does not exist in the directory
	_, err = cache.Search("missing")
	require.Error(t, err)

	// Ambiguous search results
	_, err = cache.Search("leonardo")
	require.Error(t, err)

	// Test concurrent Search calls
	t.Run("search", func(t *testing.T) {
		tests := []struct {
			name string
			peer string
		}{
			{"leonardo", "leonardo da vinci"},
			{"leonardo2", "leonardo da vinci"},
			{"donatello", "donatello"},
			{"donatello2", "donatello"},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				p, err := cache.Search(tt.peer)
				require.NoError(t, err)
				require.NotNil(t, p)
				// require.Equal(t, apimock.RemotePeers[tt.peer].ID, p.Info().ID)
			})
		}
	})

	// Cache should contain the two peers
	// leonardo, err := cache.Get("leonardo da vinci")
	// require.NoError(t, err)
	// expected := apimock.RemotePeers["leonardo da vinci"]
	// require.Equal(t, expected.ID, leonardo.Info().ID)
	// require.Equal(t, expected.RegisteredDirectory, leonardo.Info().RegisteredDirectory)
	// require.Equal(t, expected.CommonName, leonardo.Info().CommonName)
	// require.Equal(t, expected.Endpoint, leonardo.Info().Endpoint)

	// donatello, err := cache.Get("donatello")
	// require.NoError(t, err)
	// expected = apimock.RemotePeers["donatello"]
	// require.Equal(t, expected.ID, donatello.Info().ID)
	// require.Equal(t, expected.RegisteredDirectory, donatello.Info().RegisteredDirectory)
	// require.Equal(t, expected.CommonName, donatello.Info().CommonName)
	// require.Equal(t, expected.Endpoint, donatello.Info().Endpoint)
}

// Helper function to create a new Peers manager (e.g. cached peers) connected to a mock
// directory service for testing interactions with the directory service and TRISA network.
func makePeersCache() (cache *peers.Peers, mgds *gdsmock.GDS, err error) {
	// Load "client" certificates to initialize the Peers manager. It doesn't really
	// matter if the remote uses client or server or the mocked Peers cache does, they
	// just have to load a different certificate and private key than the other.
	certs, pool, err := loadCertificates("testdata/client.pem")
	if err != nil {
		return nil, nil, err
	}

	// Create the peeers cache with the configured credentials and a mock GDS
	cache = peers.New(certs, pool, "bufconn")
	mgds = gdsmock.New(nil)

	// Connect the peers cache to the mock GDS for testing purposes
	cache.Connect(
		grpc.WithContextDialer(mgds.Channel().Dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	return cache, mgds, nil
}

// Helper function to load certificates from disk
func loadCertificates(path string) (certs *trust.Provider, pool trust.ProviderPool, err error) {
	var sz *trust.Serializer
	if sz, err = trust.NewSerializer(false); err != nil {
		return nil, nil, err
	}

	if certs, err = sz.ReadFile(path); err != nil {
		return nil, nil, err
	}

	if pool, err = sz.ReadPoolFile(path); err != nil {
		return nil, nil, err
	}

	return certs, pool, nil
}
