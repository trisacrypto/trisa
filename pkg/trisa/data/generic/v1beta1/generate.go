package generic

//go:generate protoc -I=../../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/trisa/pkg/trisa/data/generic/v1beta1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/trisa/pkg/trisa/data/generic/v1beta1 trisa/data/generic/v1beta1/transaction.proto
