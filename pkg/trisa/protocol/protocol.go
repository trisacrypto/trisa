package protocol

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/trisacrypto/trisa/pkg/trisa/crypto/aesgcm"
	pb "github.com/trisacrypto/trisa/proto/trisa/protocol/v1alpha1"
)

func DecodeTransaction(ctx context.Context, t *pb.Transaction) (*pb.TransactionData, error) {

	plain, err := aesgcm.Decrypt(t.Transaction, t.Hmac, t.EncryptionKey)
	if err != nil {
		return nil, err
	}

	td := &pb.TransactionData{}
	if err := proto.Unmarshal(plain, td); err != nil {
		return nil, err
	}
	return td, nil
}

func EncodeTransactionData(ctx context.Context, id string, td *pb.TransactionData) (*pb.Transaction, error) {
	serialized, err := proto.Marshal(td)
	if err != nil {
		return nil, err
	}

	cipherText, key, sig, _, err := aesgcm.Encrypt(serialized)
	if err != nil {
		return nil, err
	}

	t := &pb.Transaction{
		Id:                  id,
		Transaction:         cipherText,
		EncryptionKey:       key,
		EncryptionAlgorithm: "AES256_GCM",
		Hmac:                sig,
		HmacSecret:          key,
		HmacAlgorithm:       "HMAC_SHA256",
	}
	return t, nil
}
