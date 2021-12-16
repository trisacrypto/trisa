package mock

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockNetworkClient struct{}

func (m *MockNetworkClient) KeyExchange(ctx context.Context, in *api.SigningKey, opts ...grpc.CallOption) (*api.SigningKey, error) {
	var (
		private *rsa.PrivateKey
		data    []byte
		err     error
	)
	if private, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to generate key pair: %v", err))
	}
	if data, err = x509.MarshalPKIXPublicKey(&private.PublicKey); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal public key: %v", err))
	}

	return &api.SigningKey{
		Data: data,
	}, nil
}

func (m *MockNetworkClient) Transfer(ctx context.Context, in *api.SecureEnvelope, opts ...grpc.CallOption) (*api.SecureEnvelope, error) {
	return nil, nil
}

func (m *MockNetworkClient) TransferStream(ctx context.Context, opts ...grpc.CallOption) (api.TRISANetwork_TransferStreamClient, error) {
	return nil, nil
}

func (m *MockNetworkClient) ConfirmAddress(ctx context.Context, in *api.Address, opts ...grpc.CallOption) (*api.AddressConfirmation, error) {
	return nil, nil
}
