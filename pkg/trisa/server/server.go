package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"

	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"github.com/trisacrypto/trisa/pkg/trisa/handler"
	"github.com/trisacrypto/trisa/pkg/trisa/protocol"
	pb "github.com/trisacrypto/trisa/proto/trisa/protocol/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Handler interface {
	HandleRequest(ctx context.Context, id string, td *pb.TransactionData) (*pb.TransactionData, error)
}

func New(h Handler, c tls.Certificate, cp *x509.CertPool) *Server {
	return &Server{
		handler:  h,
		cert:     c,
		certPool: cp,
		streams:  make(map[string]pb.TrisaPeer2Peer_TransactionStreamClient),
	}
}

type Server struct {
	handler  Handler
	cert     tls.Certificate
	certPool *x509.CertPool

	streams map[string]pb.TrisaPeer2Peer_TransactionStreamClient
}

func (s *Server) getClient(target string) (pb.TrisaPeer2Peer_TransactionStreamClient, error) {

	if stream, found := s.streams[target]; found {
		return stream, nil
	}

	u, _ := url.Parse(target)

	tls := credentials.NewTLS(&tls.Config{
		ServerName:   u.Host,
		Certificates: []tls.Certificate{s.cert},
		RootCAs:      s.certPool,
	})

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(tls))
	if err != nil {
		return nil, err
	}

	client := pb.NewTrisaPeer2PeerClient(conn)
	stream, err := client.TransactionStream(context.Background())
	if err != nil {
		return nil, err
	}

	s.streams[target] = stream
	return stream, nil
}

func (s *Server) SendRequest(ctx context.Context, target, id string, td *pb.TransactionData) error {

	ctx = handler.WithClientSide(ctx)

	t, err := protocol.EncodeTransactionData(ctx, id, td)
	if err != nil {
		return err
	}

	stream, err := s.getClient(target)

	if err != nil {
		return err
	}

	if err := stream.Send(t); err != nil {
		return err
	}

	// Extract identity
	identityType, _ := ptypes.AnyMessageName(td.Identity)
	var identityData ptypes.DynamicAny
	ptypes.UnmarshalAny(td.Identity, &identityData)

	// Extract network information
	networkType, _ := ptypes.AnyMessageName(td.Data)
	var networkData ptypes.DynamicAny
	ptypes.UnmarshalAny(td.Data, &networkData)

	log.WithFields(log.Fields{
		"identity-type": identityType,
		"network-type":  networkType,
		"identity":      fmt.Sprintf("%v", identityData),
		"network":       fmt.Sprintf("%v", networkData),
	}).Infof("sent transaction %s to %v", id, target)

	resp, err := stream.Recv()
	if err == io.EOF {
		return fmt.Errorf("premature stream exit")
	}
	if err != nil {
		return fmt.Errorf("receive stream error: %v", err)
	}

	_, err = s.handle(ctx, resp)
	if err != nil && err.Error() != "EOL" {
		return fmt.Errorf("response stream error: %v", err)
	}

	return nil
}

func (s *Server) TransactionStream(srv pb.TrisaPeer2Peer_TransactionStreamServer) error {

	ctx := srv.Context()

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			log.Info("stream exit")
			return nil
		}
		if err != nil {
			log.Warnf("receive stream error: %v", err)
			continue
		}

		resp, err := s.handle(ctx, req)
		if err != nil && err.Error() == "EOL" {
			continue
		}
		if err != nil {
			log.Warnf("response stream error: %v", err)
		}

		if err := srv.Send(resp); err != nil {
			log.Warnf("send stream error: %v", err)
		}
	}
	return nil
}

func (s *Server) handle(ctx context.Context, req *pb.Transaction) (*pb.Transaction, error) {

	log.WithFields(log.Fields{
		"direction": "incoming",
		"enc_blob":  req.Transaction,
		"enc_algo":  req.EncryptionAlgorithm,
		"hmac":      req.Hmac,
		"hmac_algo": req.HmacAlgorithm,
	}).Infof("protocol envelope for incomingtransaction %s", req.Id)

	if req.Id == "" {
		return nil, fmt.Errorf("empty transaction identifier")
	}

	reqTransactionData, err := protocol.DecodeTransaction(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decode request: %v", err)
	}

	resTransactionData, err := s.handler.HandleRequest(ctx, req.Id, reqTransactionData)
	if err != nil && err.Error() == "EOL" {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("transaction request handler request: %s", err)
	}

	res, err := protocol.EncodeTransactionData(ctx, req.Id, resTransactionData)
	if err != nil {
		return nil, fmt.Errorf("encode response: %v", err)
	}

	log.WithFields(log.Fields{
		"direction": "outgoing",
		"enc_blob":  res.Transaction,
		"enc_algo":  res.EncryptionAlgorithm,
		"hmac":      res.Hmac,
		"hmac_algo": res.HmacAlgorithm,
	}).Infof("protocol envelope for incomingtransaction %s", res.Id)

	return res, nil
}
