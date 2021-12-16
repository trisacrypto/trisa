package mock

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	gds "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type lookupInfo struct {
	ID                  string
	RegisteredDirectory string
	CommonName          string
	Endpoint            string
	ValidCertificate    bool
}

var RemotePeers = map[string]*lookupInfo{
	"leonardo": {
		ID:                  "1",
		RegisteredDirectory: "testdirectory.org",
		CommonName:          "leonardo",
		Endpoint:            "https://leonardo.trisatest.net:443",
		ValidCertificate:    true,
	},
	"donatello": {
		ID:                  "2",
		RegisteredDirectory: "testdirectory.org",
		CommonName:          "donatello",
		Endpoint:            "https://donatello.trisatest.net:443",
	},
}

type MockDirectoryClient struct{}

func (m *MockDirectoryClient) Lookup(ctx context.Context, in *gds.LookupRequest, opts ...grpc.CallOption) (*gds.LookupReply, error) {
	var (
		info *lookupInfo
		data []byte
		ok   bool
		err  error
	)
	if info, ok = RemotePeers[in.CommonName]; !ok {
		return nil, status.Error(codes.NotFound, "peer not found")
	}

	// Simulate different key types for error testing, currently only RSA is supported
	if info.ValidCertificate {
		var private *rsa.PrivateKey
		if private, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to generate key pair: %v", err))
		}
		if data, err = x509.MarshalPKIXPublicKey(&private.PublicKey); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal public key: %v", err))
		}
	} else {
		var private *ecdsa.PrivateKey
		if private, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to generate ecdsa key pair: %v", err))
		}
		if data, err = x509.MarshalPKIXPublicKey(&private.PublicKey); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to marshal public key: %v", err))
		}
	}
	return &gds.LookupReply{
		Id:                  info.ID,
		RegisteredDirectory: info.RegisteredDirectory,
		CommonName:          info.CommonName,
		Endpoint:            info.Endpoint,
		SigningCertificate:  &models.Certificate{Data: data},
	}, nil
}

func (m *MockDirectoryClient) Search(ctx context.Context, in *gds.SearchRequest, opts ...grpc.CallOption) (*gds.SearchReply, error) {
	return nil, nil
}

func (m *MockDirectoryClient) Register(ctx context.Context, in *gds.RegisterRequest, opts ...grpc.CallOption) (*gds.RegisterReply, error) {
	return nil, nil
}

func (m *MockDirectoryClient) RegisterReply(ctx context.Context, in *gds.RegisterRequest, opts ...grpc.CallOption) (*gds.RegisterReply, error) {
	return nil, nil
}

func (m *MockDirectoryClient) VerifyContact(ctx context.Context, in *gds.VerifyContactRequest, opts ...grpc.CallOption) (*gds.VerifyContactReply, error) {
	return nil, nil
}

func (m *MockDirectoryClient) Verification(ctx context.Context, in *gds.VerificationRequest, opts ...grpc.CallOption) (*gds.VerificationReply, error) {
	return nil, nil
}

func (m *MockDirectoryClient) Status(ctx context.Context, in *gds.HealthCheck, opts ...grpc.CallOption) (*gds.ServiceState, error) {
	return nil, nil
}
