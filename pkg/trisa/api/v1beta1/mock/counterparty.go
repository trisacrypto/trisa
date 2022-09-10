package mock

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/trisacrypto/trisa/pkg/bufconn"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"github.com/trisacrypto/trisa/pkg/trisa/mtls"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	TransferRPC       = "trisa.api.v1beta1.TRISANetwork/Transfer"
	TransferStreamRPC = "trisa.api.v1beta1.TRISANetwork/TransferStream"
	KeyExchangeRPC    = "trisa.api.v1beta1.TRISANetwork/KeyExchange"
	ConfirmAddressRPC = "trisa.api.v1beta1.TRISANetwork/ConfirmAddress"
	StatusRPC         = "trisa.api.v1beta1.TRISAHealth/Status"
)

// New creates a new mock RemotePeer. If bufnet is nil, one is created for the user.
func New(bufnet *bufconn.Listener) *RemotePeer {
	if bufnet == nil {
		bufnet = bufconn.New()
	}

	remote := &RemotePeer{
		bufnet: bufnet,
		srv:    grpc.NewServer(),
		Calls:  make(map[string]int),
	}

	api.RegisterTRISANetworkServer(remote.srv, remote)
	api.RegisterTRISAHealthServer(remote.srv, remote)
	go remote.srv.Serve(remote.bufnet.Sock())
	return remote
}

// NewAuth creates a new mock RemotePeer that enforces mTLS credentials over bufconn.
func NewAuth(bufnet *bufconn.Listener, certs *trust.Provider, pool trust.ProviderPool) (remote *RemotePeer, err error) {
	if bufnet == nil {
		bufnet = bufconn.New()
	}

	var creds grpc.ServerOption
	if creds, err = mtls.ServerCreds(certs, pool); err != nil {
		return nil, err
	}

	remote = &RemotePeer{
		bufnet: bufnet,
		srv:    grpc.NewServer(creds),
		Calls:  make(map[string]int),
	}

	api.RegisterTRISANetworkServer(remote.srv, remote)
	api.RegisterTRISAHealthServer(remote.srv, remote)
	go remote.srv.Serve(remote.bufnet.Sock())
	return remote, nil
}

// RemotePeer implements a mock gRPC server for testing Peer client connections. The
// desired response of the remote peer can be set by external callers using the OnRPC
// functions or the WithFixture or WithError functions. The Calls map can be used to
// count the number of times the remote peer RPC was called.
type RemotePeer struct {
	sync.Mutex
	api.UnimplementedTRISAHealthServer
	api.UnimplementedTRISANetworkServer
	bufnet           *bufconn.Listener
	srv              *grpc.Server
	Calls            map[string]int
	OnTransfer       func(context.Context, *api.SecureEnvelope) (*api.SecureEnvelope, error)
	OnTransferStream func(api.TRISANetwork_TransferStreamServer) error
	OnKeyExchange    func(context.Context, *api.SigningKey) (*api.SigningKey, error)
	OnConfirmAddress func(context.Context, *api.Address) (*api.AddressConfirmation, error)
	OnStatus         func(context.Context, *api.HealthCheck) (*api.ServiceState, error)
}

func (s *RemotePeer) Channel() *bufconn.Listener {
	return s.bufnet
}

func (s *RemotePeer) Shutdown() {
	s.srv.GracefulStop()
	s.bufnet.Close()
}

func (s *RemotePeer) Reset() {
	for key := range s.Calls {
		s.Calls[key] = 0
	}

	s.OnTransfer = nil
	s.OnTransferStream = nil
	s.OnKeyExchange = nil
	s.OnConfirmAddress = nil
	s.OnStatus = nil
}

// UseFixture loadsa a JSON fixture from disk (usually in a testdata folder) to use as
// the protocol buffer response to the specified RPC, simplifying handler mocking.
func (s *RemotePeer) UseFixture(rpc, path string) (err error) {
	var data []byte
	if data, err = os.ReadFile(path); err != nil {
		return fmt.Errorf("could not read fixture: %v", err)
	}

	jsonpb := &protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	switch rpc {
	case TransferRPC:
		out := &api.SecureEnvelope{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnTransfer = func(context.Context, *api.SecureEnvelope) (*api.SecureEnvelope, error) {
			return out, nil
		}
	case TransferStreamRPC:
		return errors.New("cannot use fixture for a streaming RPC")
	case KeyExchangeRPC:
		out := &api.SigningKey{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnKeyExchange = func(context.Context, *api.SigningKey) (*api.SigningKey, error) {
			return out, nil
		}
	case ConfirmAddressRPC:
		out := &api.AddressConfirmation{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnConfirmAddress = func(context.Context, *api.Address) (*api.AddressConfirmation, error) {
			return out, nil
		}
	case StatusRPC:
		out := &api.ServiceState{}
		if err = jsonpb.Unmarshal(data, out); err != nil {
			return fmt.Errorf("could not unmarshal json into %T: %v", out, err)
		}
		s.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			return out, nil
		}
	default:
		return fmt.Errorf("unknown RPC %q", rpc)
	}

	return nil
}

// UseError allows you to specify a gRPC status error to return from the specified RPC.
func (s *RemotePeer) UseError(rpc string, code codes.Code, msg string) error {
	switch rpc {
	case TransferRPC:
		s.OnTransfer = func(context.Context, *api.SecureEnvelope) (*api.SecureEnvelope, error) {
			return nil, status.Error(code, msg)
		}
	case TransferStreamRPC:
		s.OnTransferStream = func(api.TRISANetwork_TransferStreamServer) error {
			return status.Error(code, msg)
		}
	case KeyExchangeRPC:
		s.OnKeyExchange = func(context.Context, *api.SigningKey) (*api.SigningKey, error) {
			return nil, status.Error(code, msg)
		}
	case ConfirmAddressRPC:
		s.OnConfirmAddress = func(context.Context, *api.Address) (*api.AddressConfirmation, error) {
			return nil, status.Error(code, msg)
		}
	case StatusRPC:
		s.OnStatus = func(context.Context, *api.HealthCheck) (*api.ServiceState, error) {
			return nil, status.Error(code, msg)
		}
	default:
		return fmt.Errorf("unknown RPC %q", rpc)
	}
	return nil
}

func (s *RemotePeer) Transfer(ctx context.Context, in *api.SecureEnvelope) (*api.SecureEnvelope, error) {
	s.IncrementCalls(TransferRPC)
	return s.OnTransfer(ctx, in)
}

func (s *RemotePeer) TransferStream(stream api.TRISANetwork_TransferStreamServer) error {
	s.IncrementCalls(TransferStreamRPC)
	return s.OnTransferStream(stream)
}

func (s *RemotePeer) KeyExchange(ctx context.Context, in *api.SigningKey) (*api.SigningKey, error) {
	s.IncrementCalls(KeyExchangeRPC)
	return s.OnKeyExchange(ctx, in)
}

func (s *RemotePeer) ConfirmAddress(ctx context.Context, in *api.Address) (*api.AddressConfirmation, error) {
	s.IncrementCalls(ConfirmAddressRPC)
	return s.OnConfirmAddress(ctx, in)
}

func (s *RemotePeer) Status(ctx context.Context, in *api.HealthCheck) (*api.ServiceState, error) {
	s.IncrementCalls(StatusRPC)
	return s.OnStatus(ctx, in)
}

func (s *RemotePeer) IncrementCalls(rpc string) {
	s.Lock()
	s.Calls[rpc]++
	s.Unlock()
}
