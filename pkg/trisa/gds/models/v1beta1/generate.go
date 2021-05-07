package models

//go:generate protoc -I=../../../../../proto --go_out=. --go_opt=module=github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1 --go-grpc_out=. --go-grpc_opt=module=github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1 trisa/gds/models/v1beta1/ca.proto trisa/gds/models/v1beta1/models.proto
