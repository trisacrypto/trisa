package peers_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	apimock "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1/mock"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	gdsmock "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1/mock"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trisa/peers"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Test that Add correctly adds peers to the Peers cache.
func TestAdd(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Common name is required to add a peer to the cache
	err = cache.Add(&peers.PeerInfo{})
	require.EqualError(t, err, "common name is required for all peers")

	// Generate a random key for some of our fixtures.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Test adding peers concurrently; the leonardo peer should be updated with
	// consecutive updates to
	t.Run("addTests", func(t *testing.T) {
		tests := []struct {
			name string
			info *peers.PeerInfo
		}{
			{"add-id-only", &peers.PeerInfo{
				CommonName: "leonardo.trisatest.net",
				ID:         "19d84515-007a-48cc-9efd-b153a263e77c",
			}},
			{"add-registered-directory-only", &peers.PeerInfo{
				CommonName:          "leonardo.trisatest.net",
				RegisteredDirectory: "testdirectory.org",
			}},
			{"add-endpoint-only", &peers.PeerInfo{
				CommonName: "leonardo.trisatest.net",
				Endpoint:   "leonardo.trisatest.net:443",
			}},
			{"add-signing-key-only", &peers.PeerInfo{
				CommonName: "leonardo.trisatest.net",
				SigningKey: &privateKey.PublicKey,
			}},
			{"add-new-peer", &peers.PeerInfo{
				CommonName:          "donatello.trisatest.net",
				ID:                  "b19c9ebd-82f5-4bda-91ef-226e3ecee4b8",
				RegisteredDirectory: "testdirectory.org",
				Endpoint:            "donatello.trisatest.net:443",
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
	leonardo, err := cache.Get("leonardo.trisatest.net")
	require.NoError(t, err)
	require.Equal(t, "leonardo.trisatest.net", leonardo.Info().CommonName)
	require.Equal(t, "19d84515-007a-48cc-9efd-b153a263e77c", leonardo.Info().ID)
	require.Equal(t, "testdirectory.org", leonardo.Info().RegisteredDirectory)
	require.Equal(t, "leonardo.trisatest.net:443", leonardo.Info().Endpoint)
	require.Equal(t, &privateKey.PublicKey, leonardo.Info().SigningKey)

	donatello, err := cache.Get("donatello.trisatest.net")
	require.NoError(t, err)
	require.Equal(t, "donatello.trisatest.net", donatello.Info().CommonName)
	require.Equal(t, "b19c9ebd-82f5-4bda-91ef-226e3ecee4b8", donatello.Info().ID)
	require.Equal(t, "testdirectory.org", donatello.Info().RegisteredDirectory)
	require.Equal(t, "donatello.trisatest.net:443", donatello.Info().Endpoint)
}

// Test that FromContext returns the correct Peer given the connection context.
func TestFromContext(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Make assertions about what is being looked up in the GDS
	mgds.OnLookup = func(_ context.Context, in *gds.LookupRequest) (out *gds.LookupReply, err error) {
		// Assert that the expected common name is being looked up
		require.Equal(t, "server.trisa.dev", in.CommonName, "unexpected common name in lookup request")
		require.Empty(t, in.Id, "unexpected id in lookup request")
		require.Empty(t, in.RegisteredDirectory, "unexpected registered directory in lookup request")

		return &gds.LookupReply{
			Id:                  "0960c00e-68a7-4606-9d0f-ff8537186d34",
			RegisteredDirectory: "localhost",
			CommonName:          "server.trisa.dev",
			Endpoint:            "server.trisa.dev:4000",
			Name:                "Testing VASP",
			Country:             "US",
			VerifiedOn:          "2022-05-10T22:29:55Z",
		}, nil
	}

	// Create an mTLS connection to test the context over bufconn
	certs, pool, err := loadCertificates("testdata/server.pem")
	require.NoError(t, err, "could not load server certificate fixtures")

	remote, err := apimock.NewAuth(nil, certs, pool)
	require.NoError(t, err, "could not create remote universal mock")
	defer remote.Shutdown()

	// Lookup the peer from the directory service and connect it to the remote
	peer, err := cache.Lookup("server.trisa.dev")
	require.NoError(t, err, "could not lookup peer from directory service")

	opts := make([]grpc.DialOption, 0, 2)
	opts = append(opts, grpc.WithContextDialer(remote.Channel().Dialer))

	clientCerts, clientPool, err := loadCertificates("testdata/client.pem")
	require.NoError(t, err, "could not load client certificate fixture")
	creds, err := mtls.ClientCreds(peer.Info().Endpoint, clientCerts, clientPool)
	require.NoError(t, err, "could not create mtls client credentials")
	opts = append(opts, creds)

	err = peer.Connect(opts...)
	require.NoError(t, err, "could not connect to remote peer with mtls credentials")

	// Setup to get the context from the remote dialer
	remote.OnTransfer = func(ctx context.Context, _ *api.SecureEnvelope) (*api.SecureEnvelope, error) {
		// Ok, after all that work above we finally have an actual gRPC context with mTLS info
		peer, err := cache.FromContext(ctx)
		require.NoError(t, err, "could not lookup peer from context")

		info := peer.Info()
		require.Equal(t, "client.trisa.dev", info.CommonName, "unknown common name")

		// Don't return anything
		return &api.SecureEnvelope{}, nil
	}

	// Make the request with the client to finish the tests
	_, err = peer.Transfer(&api.SecureEnvelope{})
	require.NoError(t, err, "could not make transfer to initiate from context tests")

	// TODO: Test with different certificates to make sure the auth info parsing in FromContext is correct
}

// Test that the Lookup function returns the correct remote peer given the common name.
func TestLookup(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Handle the case where the GDS returns an error
	mgds.UseError(gdsmock.LookupRPC, codes.NotFound, "could not find peer with that common name")
	peer, err := cache.Lookup("unknown")
	require.EqualError(t, err, "rpc error: code = NotFound desc = could not find peer with that common name")
	require.Nil(t, peer, "peer should be nil when an error is returned")

	// Handle the case where the GDS returns an error in the lookup reply
	mgds.OnLookup = func(context.Context, *gds.LookupRequest) (*gds.LookupReply, error) {
		return &gds.LookupReply{
			Error: &gds.Error{
				Code:    99,
				Message: "the GDS really shouldn't be returning these errors",
			},
		}, nil
	}

	peer, err = cache.Lookup("unknown")
	require.EqualError(t, err, "[99] the GDS really shouldn't be returning these errors")
	require.Nil(t, peer, "peer should be nil when an error is returned")

	// Lookup should not error if GDS returns a nil identity and nil signing certificate
	require.NoError(t, mgds.UseFixture(gdsmock.LookupRPC, "testdata/leonardo.trisa.dev.pb.json"))
	peer, err = cache.Lookup("leonardo.trisa.dev")
	require.NoError(t, err, "could not lookup peer from directory service")
	require.Nil(t, peer.SigningKey(), "signing key should be nil when no certificate is returned")

	// Configure a reply fixture for the other lookup test cases
	reply := &gds.LookupReply{}
	require.NoError(t, loadGRPCFixture("testdata/leonardo.trisa.dev.pb.json", reply))
	mgds.OnLookup = func(context.Context, *gds.LookupRequest) (*gds.LookupReply, error) {
		return reply, nil
	}

	// Generate keys to distinguish between identity and signing certificates
	origCert, origKey, err := generateCertificate()
	require.NoError(t, err, "could not generate certificate")
	identityCert, _, err := generateCertificate()
	require.NoError(t, err, "could not generate certificate")
	signingCert, signingKey, err := generateCertificate()
	require.NoError(t, err, "could not generate certificate")

	// Lookup should store the identity key if the identity certificate is available on the lookup reply
	reply.CommonName = "identity"
	reply.IdentityCertificate = origCert
	peer, err = cache.Lookup("identity")
	require.NoError(t, err, "could not lookup peer from directory service")
	require.Equal(t, origKey, peer.SigningKey(), "identity key should be stored")

	// Lookup should not overwrite an existing key on the peer
	reply.IdentityCertificate = identityCert
	peer, err = cache.Lookup("identity")
	require.NoError(t, err, "could not lookup peer from directory service")
	require.Equal(t, origKey, peer.SigningKey(), "identity key should not be overwritten")

	// Lookup should store the signing key if the signing certificate is available on the lookup reply
	reply.CommonName = "signing"
	reply.IdentityCertificate = nil
	reply.SigningCertificate = signingCert
	peer, err = cache.Lookup("signing")
	require.NoError(t, err, "could not lookup peer from directory service")
	require.Equal(t, signingKey, peer.SigningKey(), "signing key should be stored")

	// Lookup should prefer the signing certificate over the identity certificate if both are available on the lookup reply
	reply.CommonName = "both"
	reply.IdentityCertificate = identityCert
	reply.SigningCertificate = signingCert
	peer, err = cache.Lookup("both")
	require.NoError(t, err, "could not lookup peer from directory service")
	require.Equal(t, signingKey, peer.SigningKey(), "signing key should be stored")

	// Handle the case where the GDS returns valid responses
	mgds.Reset()
	mgds.OnLookup = func(_ context.Context, in *gds.LookupRequest) (out *gds.LookupReply, err error) {
		out = &gds.LookupReply{}
		switch in.CommonName {
		case "leonardo.trisa.dev":
			if err = loadGRPCFixture("testdata/leonardo.trisa.dev.pb.json", out); err != nil {
				return nil, err
			}
		case "donatello.example.com":
			if err = loadGRPCFixture("testdata/donatello.example.com.pb.json", out); err != nil {
				return nil, err
			}
		default:
			return nil, status.Error(codes.NotFound, "unknown TRISA counterparty")
		}
		return out, nil
	}

	// Test concurrent Lookup calls
	t.Run("lookup", func(t *testing.T) {
		tests := []struct {
			name string
			peer string
		}{
			{"lookup-leonardo", "leonardo.trisa.dev"},
			{"lookup-donatello", "donatello.example.com"},
		}
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				p, err := cache.Lookup(tt.peer)
				require.NoError(t, err)
				require.NotNil(t, p)
			})
		}
	})

	require.Equal(t, 0, mgds.Calls[gdsmock.SearchRPC])
	require.Equal(t, 2, mgds.Calls[gdsmock.LookupRPC])
	require.Equal(t, 0, mgds.Calls[gdsmock.StatusRPC])

	// Cache should contain the two peers
	leonardo, err := cache.Get("leonardo.trisa.dev")
	require.NoError(t, err)
	require.Equal(t, "19d84515-007a-48cc-9efd-b153a263e77c", leonardo.Info().ID)
	require.Equal(t, "testdirectory.org", leonardo.Info().RegisteredDirectory)
	require.Equal(t, "leonardo.trisa.dev", leonardo.Info().CommonName)
	require.Equal(t, "leonardo.trisa.dev:8000", leonardo.Info().Endpoint)

	donatello, err := cache.Get("donatello.example.com")
	require.NoError(t, err)
	require.Equal(t, "b19c9ebd-82f5-4bda-91ef-226e3ecee4b8", donatello.Info().ID)
	require.Equal(t, "testdirectory.org", donatello.Info().RegisteredDirectory)
	require.Equal(t, "donatello.example.com", donatello.Info().CommonName)
	require.Equal(t, "donatello.example.com:443", donatello.Info().Endpoint)
}

// Test that the Search function returns the matching remote peer given the name.
func TestSearch(t *testing.T) {
	// Create a mocked peers cache connected to a mock directory
	cache, mgds, err := makePeersCache()
	require.NoError(t, err, "could not create mocked peers cache")
	defer mgds.Shutdown()

	// Handle the case where GDS returns an error
	mgds.UseError(gdsmock.SearchRPC, codes.NotFound, "the search terms you provided were not found")
	_, err = cache.Search("missing")
	require.EqualError(t, err, "rpc error: code = NotFound desc = the search terms you provided were not found")

	// Handle the case where the GDS returns an error in the lookup reply
	mgds.OnSearch = func(context.Context, *gds.SearchRequest) (*gds.SearchReply, error) {
		return &gds.SearchReply{
			Error: &gds.Error{
				Code:    99,
				Message: "the GDS really shouldn't be returning these errors",
			},
		}, nil
	}

	_, err = cache.Search("missing")
	require.EqualError(t, err, "[99] the GDS really shouldn't be returning these errors")

	// Handle the case where GDS returns no results in the search reply
	mgds.OnSearch = func(context.Context, *gds.SearchRequest) (*gds.SearchReply, error) {
		return &gds.SearchReply{
			Results: make([]*gds.SearchReply_Result, 0),
		}, nil
	}

	_, err = cache.Search("Da Vinci Digital Exchange")
	require.EqualError(t, err, "could not find peer named \"Da Vinci Digital Exchange\"")

	err = mgds.UseFixture(gdsmock.SearchRPC, "testdata/gds_search_reply.pb.json")
	require.NoError(t, err, "could not load multiple results fixture")

	_, err = cache.Search("Da Vinci Digital Exchange")
	require.EqualError(t, err, "too many results returned for \"Da Vinci Digital Exchange\"")

	// Have the mock GDS respond correctly based on the input
	mgds.Reset()
	mgds.OnSearch = func(_ context.Context, in *gds.SearchRequest) (out *gds.SearchReply, err error) {
		out = &gds.SearchReply{}
		if err = loadGRPCFixture("testdata/gds_search_reply.pb.json", out); err != nil {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		lookupID := map[string]string{
			"Da Vinci Digital Exchange": "19d84515-007a-48cc-9efd-b153a263e77c",
			"Brooklyn BitMining Ltd":    "b19c9ebd-82f5-4bda-91ef-226e3ecee4b8",
		}[in.Name[0]]

		results := make([]*gds.SearchReply_Result, 0, 1)
		for _, result := range out.Results {
			if result.Id == lookupID {
				results = append(results, result)
			}
		}

		return &gds.SearchReply{
			Results: results,
		}, nil
	}

	// Test concurrent Search calls populating the cache
	t.Run("search", func(t *testing.T) {
		tests := []struct {
			name string
			peer string
		}{
			{"search-leonardo", "Da Vinci Digital Exchange"},
			{"search-donatello", "Brooklyn BitMining Ltd"},
		}
		for _, tt := range tests {
			tc := tt
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				p, err := cache.Search(tc.peer)
				require.NoError(t, err)
				require.NotNil(t, p)
			})
		}
	})

	require.Equal(t, 2, mgds.Calls[gdsmock.SearchRPC])
	require.Equal(t, 0, mgds.Calls[gdsmock.LookupRPC])
	require.Equal(t, 0, mgds.Calls[gdsmock.StatusRPC])

	// Cache should contain the two peers
	leonardo, err := cache.Get("leonardo.trisa.dev")
	require.NoError(t, err)
	require.Equal(t, "19d84515-007a-48cc-9efd-b153a263e77c", leonardo.Info().ID)
	require.Equal(t, "testdirectory.org", leonardo.Info().RegisteredDirectory)
	require.Equal(t, "leonardo.trisa.dev", leonardo.Info().CommonName)
	require.Equal(t, "leonardo.trisa.dev:8000", leonardo.Info().Endpoint)

	donatello, err := cache.Get("donatello.example.com")
	require.NoError(t, err)
	require.Equal(t, "b19c9ebd-82f5-4bda-91ef-226e3ecee4b8", donatello.Info().ID)
	require.Equal(t, "testdirectory.org", donatello.Info().RegisteredDirectory)
	require.Equal(t, "donatello.example.com", donatello.Info().CommonName)
	require.Equal(t, "donatello.example.com:443", donatello.Info().Endpoint)
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

// Helper function to generate a certificate with a random key
func generateCertificate() (cert *models.Certificate, key *rsa.PublicKey, err error) {
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, nil, err
	}

	key = &privateKey.PublicKey
	cert = &models.Certificate{}
	if cert.Data, err = x509.MarshalPKIXPublicKey(key); err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

// Helper function to load gRPC fixtures from disks, errors will be status errors.
func loadGRPCFixture(path string, v proto.Message) (err error) {
	pbjson := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not read fixture: %s", err)
	}

	if err = pbjson.Unmarshal(data, v); err != nil {
		return status.Errorf(codes.FailedPrecondition, "could not unmarshal fixture: %s", err)
	}

	return nil
}
