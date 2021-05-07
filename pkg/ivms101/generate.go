package ivms101

//go:generate protoc -I=../../proto --go_out=. --go_opt=module=github.com/trisacrypto/trisa/pkg/ivms101 ivms101/enum.proto ivms101/ivms101.proto ivms101/identity.proto
