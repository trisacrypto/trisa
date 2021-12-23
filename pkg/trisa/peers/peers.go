package peers

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1/mock"
	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// Peers manages TRISA network connections to send requests to other TRISA nodes.
type Peers struct {
	sync.RWMutex
	certs        *trust.Provider
	pool         trust.ProviderPool
	peers        map[string]*Peer
	directoryURL string
	directory    gds.TRISADirectoryClient
}

// New creates a new Peers cache to look up peers from context or by endpoint.
func New(certs *trust.Provider, pool trust.ProviderPool, directoryURL string) *Peers {
	p := &Peers{
		certs:        certs,
		pool:         pool,
		peers:        make(map[string]*Peer),
		directoryURL: directoryURL,
	}
	return p
}

// NewMock creates a mocked Peers cache that should only be used by tests. The mock
// cache contains a single peer with the common name "test-peer" and a mocked
// TRISANetworkClient which returns canned responses for the peer-to-peer network
// requests, to avoid having to run a test server. Similarly, the mock cache contains a
// mocked TRISADirectoryClient which returns canned responses for the Peers directory
// service requests.
func NewMock(certs *trust.Provider, pool trust.ProviderPool, directoryURL string) *Peers {
	const peerName = "test-peer"
	p := &Peers{
		certs:        certs,
		pool:         pool,
		peers:        make(map[string]*Peer),
		directoryURL: directoryURL,
		directory:    &mock.MockDirectoryClient{},
	}
	p.peers[peerName] = &Peer{
		parent: p,
		info: &PeerInfo{
			CommonName: peerName,
			SigningKey: &rsa.PublicKey{},
		},
		client: &mock.MockNetworkClient{},
	}
	return p
}

// Add creates or updates a peer in the peers cache with the specified info.
func (p *Peers) Add(info *PeerInfo) (err error) {
	if info.CommonName == "" {
		return errors.New("common name is required for all peers")
	}

	// Critical section for Peers
	var peer *Peer
	if peer, err = p.Get(info.CommonName); err != nil {
		return err
	}

	// Critical section for peer
	// Only update if data is available on info to avoid overwriting existing data
	peer.Lock()
	if info.ID != "" {
		peer.info.ID = info.ID
	}
	if info.RegisteredDirectory != "" {
		peer.info.RegisteredDirectory = info.RegisteredDirectory
	}
	if info.Endpoint != "" {
		peer.info.Endpoint = info.Endpoint
	}
	if info.SigningKey != nil {
		peer.info.SigningKey = info.SigningKey
	}
	peer.Unlock()
	return nil
}

// FromContext looks up the TLSInfo from the incoming gRPC connection to get the common
// name of the Peer from the certificate. If the Peer is already in the cache, it
// returns the peer information, otherwise it creates and caches the Peer info.
func (p *Peers) FromContext(ctx context.Context) (_ *Peer, err error) {
	var (
		ok         bool
		gp         *peer.Peer
		tlsAuth    credentials.TLSInfo
		commonName string
	)

	if gp, ok = peer.FromContext(ctx); !ok {
		return nil, errors.New("no peer found in context")
	}

	if tlsAuth, ok = gp.AuthInfo.(credentials.TLSInfo); !ok {
		return nil, fmt.Errorf("unexpected peer transport credentials type: %T", gp.AuthInfo)
	}

	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return nil, errors.New("could not verify peer certificate")
	}

	commonName = tlsAuth.State.VerifiedChains[0][0].Subject.CommonName
	if commonName == "" {
		return nil, errors.New("could not find common name on authenticated subject")
	}

	// Critical section
	return p.Get(commonName)
}

// Get a cached peer by common name, creating it if necessary. Getting the Peer does
// not necessarily guarantee the peer with the common name exists
func (p *Peers) Get(commonName string) (*Peer, error) {
	var (
		ok   bool
		peer *Peer
	)

	p.Lock()
	// Check if peer is already cached in memory. If not, add the new peer.
	if peer, ok = p.peers[commonName]; !ok {
		peer = &Peer{
			parent: p,
			info:   &PeerInfo{CommonName: commonName},
		}
		p.peers[commonName] = peer

		// TODO: Do a directory service lookup for the ID and registered ID
	}
	p.Unlock()
	return peer, nil
}

// Lookup uses the directory service to find the remote peer by common name.
func (p *Peers) Lookup(commonName string) (peer *Peer, err error) {
	// Lookup the peer to ensure that a peer with common name is cached.
	if peer, err = p.Get(commonName); err != nil {
		return nil, err
	}

	// Ensure we're connected to the directory service
	if err = p.Connect(); err != nil {
		return nil, err
	}

	var rep *gds.LookupReply
	req := &gds.LookupRequest{
		CommonName: commonName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if rep, err = p.directory.Lookup(ctx, req); err != nil {
		return nil, err
	}

	if rep.Error != nil {
		return nil, rep.Error
	}

	// Parse the signing certificate and create the Peer info.
	// If no signing certificates available, just ignore to do peer exchange.
	info := &PeerInfo{
		ID:                  rep.Id,
		RegisteredDirectory: rep.RegisteredDirectory,
		CommonName:          rep.CommonName,
		Endpoint:            rep.Endpoint,
	}

	var pub interface{}
	if pub, err = x509.ParsePKIXPublicKey(rep.SigningCertificate.Data); err == nil {
		var ok bool
		if info.SigningKey, ok = pub.(*rsa.PublicKey); !ok {
			info.SigningKey = nil
		}
	}

	// Update the info on the peers
	if err = p.Add(info); err != nil {
		return nil, err
	}

	return peer, nil
}

// Search uses the directory service to find a remote peer by name
func (p *Peers) Search(name string) (_ *Peer, err error) {
	// Ensure we're connected to the directory service
	if err = p.Connect(); err != nil {
		return nil, err
	}

	var rep *gds.SearchReply
	req := &gds.SearchRequest{
		Name: []string{name},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if rep, err = p.directory.Search(ctx, req); err != nil {
		return nil, err
	}

	if rep.Error != nil {
		return nil, rep.Error
	}

	if len(rep.Results) == 0 {
		return nil, fmt.Errorf("could not find peer named %q", name)
	}

	if len(rep.Results) > 1 {
		return nil, fmt.Errorf("too many results returned for %q", name)
	}

	// Create the peer info and update the cache
	info := &PeerInfo{
		ID:                  rep.Results[0].Id,
		RegisteredDirectory: rep.Results[0].RegisteredDirectory,
		CommonName:          rep.Results[0].CommonName,
		Endpoint:            rep.Results[0].Endpoint,
	}

	// Update the info on the peers
	if err = p.Add(info); err != nil {
		return nil, err
	}
	return p.Get(info.CommonName)
}

// Connect to the remote peer - thread safe.
func (p *Peers) Connect() (err error) {
	p.Lock()
	err = p.connect()
	p.Unlock()
	return err
}

// Connect to the remote peer - not thread safe.
func (p *Peers) connect() (err error) {
	// Are we already connected?
	if p.directory != nil {
		return nil
	}

	if p.directoryURL == "" {
		return errors.New("no directory service URL to dial")
	}

	opts := make([]grpc.DialOption, 0, 1)
	config := &tls.Config{}
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))

	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(p.directoryURL, opts...); err != nil {
		return err
	}

	p.directory = gds.NewTRISADirectoryClient(cc)
	return nil
}
