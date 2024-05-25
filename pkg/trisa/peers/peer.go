package peers

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"

	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"google.golang.org/grpc"
)

// Peer contains cached information about connections to other members of the TRISA
// network and facilitates directory service lookups and information exchanges.
// TODO: implement transfer stream and account confirmation endpoints.
type Peer struct {
	sync.RWMutex
	parent *Peers    // Contains common configuration for all peers
	info   *PeerInfo // NOTE: common name cannot be modified after init, see String()
	client api.TRISANetworkClient
}

// PeerInfo contains directory service information that uniquely identifies the peer.
// It is maintained separately from the Peer to allow for thread-safe reads and simpler
// marshalling and unmarshalling of JSON data about the peer.
//
// TODO: implement Marshaler and Unmarshaler to ensure signing key is base64 PEM encoded.
// TODO: allow different signing key types other than just RSA
type PeerInfo struct {
	ID                  string
	RegisteredDirectory string
	CommonName          string
	Endpoint            string
	SigningKey          *rsa.PublicKey
}

// SigningKey returns the current signing key of the remote peer, if it's available
// (otherwise returns nil). If a key exchange is underway, this method blocks until a
// key has been retrieved from the remote peer.
func (p *Peer) SigningKey() *rsa.PublicKey {
	p.RLock()
	defer p.RUnlock()
	return p.info.SigningKey
}

// UpdateSigningKey if the key exchange was initiated from a remote TRISA peer.
func (p *Peer) UpdateSigningKey(key interface{}) error {
	p.Lock()
	defer p.Unlock()

	var ok bool
	if p.info.SigningKey, ok = key.(*rsa.PublicKey); !ok {
		return fmt.Errorf("unsupported public key type %T", key)
	}
	return nil
}

// ExchangeKeys kicks of a key exchange with the remote peer. It locks to block multiple
// key exchanges from being issued and returns the key immediately if the key is already
// cached on the Peer (unless force is specified, then it will conduct a key exchange).
// This allows callers to ensure that they will get the public signing key when needed.
func (p *Peer) ExchangeKeys(force bool) (_ *rsa.PublicKey, err error) {
	// This lock causes everyone who wants the public key of the peer to wait until the
	// key exchange has been completed, reducing the number of retries overall.
	// This lock will contend with the RLock in SigningKeys() and the locking performance
	// could be improved; we err on the side of a hard lock for safety in the MVP.
	p.Lock()
	defer p.Unlock()

	// If force - then set the signing key to nil to ensure a key exchange occurs.
	if force {
		p.info.SigningKey = nil
	}

	// If we have the signing key already, just return it.
	if p.info.SigningKey != nil {
		return p.info.SigningKey, nil
	}

	var localKey *x509.Certificate
	if localKey, err = p.parent.certs.GetLeafCertificate(); err != nil {
		return nil, fmt.Errorf("invalid local signing key: %s", err)
	}

	// Create the key exchange request
	req := &api.SigningKey{
		Version:            int64(localKey.Version),
		Signature:          localKey.Signature,
		SignatureAlgorithm: localKey.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: localKey.PublicKeyAlgorithm.String(),
		NotBefore:          localKey.NotBefore.Format(time.RFC3339),
		NotAfter:           localKey.NotAfter.Format(time.RFC3339),
	}

	if req.Data, err = x509.MarshalPKIXPublicKey(localKey.PublicKey); err != nil {
		return nil, fmt.Errorf("could not marshal PKIX public key: %s", err)
	}

	// Connect to the client if not already connected
	if err = p.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var rep *api.SigningKey
	if rep, err = p.client.KeyExchange(ctx, req); err != nil {
		return nil, err
	}

	// Parse public keys from remote client
	// TODO: does this function marshal the key correctly? We should use the keys package.
	var pub interface{}
	if pub, err = x509.ParsePKIXPublicKey(rep.Data); err != nil {
		return nil, err
	}

	var ok bool
	if p.info.SigningKey, ok = pub.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("unsupported public key type %T", pub)
	}

	return p.info.SigningKey, nil
}

// Transfer sends the unary RPC request via the peer client, ensuring its connected.
func (p *Peer) Transfer(in *api.SecureEnvelope) (out *api.SecureEnvelope, err error) {
	// Thread-safe assurance that we're connected to the remote peer.
	if err = p.Connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return p.client.Transfer(ctx, in)
}

// Info returns details about the remote Peer.
func (p *Peer) Info() PeerInfo {
	p.RLock()
	defer p.RUnlock()
	return *p.info
}

// String returns the common name of the peer
func (p *Peer) String() string {
	// NOTE: common name is only modified when the Peer is created, so it is safe to
	// access the common name on the Peer as long as this rule is maintained.
	return p.info.CommonName
}

// Connect to the remote peer - thread safe.
func (p *Peer) Connect(opts ...grpc.DialOption) (err error) {
	p.Lock()
	err = p.connect(opts...)
	p.Unlock()
	return err
}

// Connect to the remote peer - not thread safe.
func (p *Peer) connect(opts ...grpc.DialOption) (err error) {
	// Are we already connected?
	if p.client != nil {
		return nil
	}

	if p.info.Endpoint == "" {
		return errors.New("peer does not have an endpoint to connect to")
	}

	if len(opts) == 0 {
		opts = make([]grpc.DialOption, 0, 1)

		var opt grpc.DialOption
		if opt, err = mtls.ClientCreds(p.info.Endpoint, p.parent.certs, p.parent.pool); err != nil {
			return err
		}

		opts = append(opts, opt)
	}

	var cc *grpc.ClientConn
	if cc, err = grpc.NewClient(p.info.Endpoint, opts...); err != nil {
		return err
	}

	p.client = api.NewTRISANetworkClient(cc)
	return nil
}
