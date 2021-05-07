package api

//go:generate protoc -I=../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1 trisa/api/v1beta1/api.proto trisa/api/v1beta1/errors.proto
