// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: trisa/api/v1beta1/api.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ServiceState_Status int32

const (
	ServiceState_UNKNOWN     ServiceState_Status = 0
	ServiceState_HEALTHY     ServiceState_Status = 1
	ServiceState_UNHEALTHY   ServiceState_Status = 2
	ServiceState_DANGER      ServiceState_Status = 3
	ServiceState_OFFLINE     ServiceState_Status = 4
	ServiceState_MAINTENANCE ServiceState_Status = 5
)

// Enum value maps for ServiceState_Status.
var (
	ServiceState_Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "HEALTHY",
		2: "UNHEALTHY",
		3: "DANGER",
		4: "OFFLINE",
		5: "MAINTENANCE",
	}
	ServiceState_Status_value = map[string]int32{
		"UNKNOWN":     0,
		"HEALTHY":     1,
		"UNHEALTHY":   2,
		"DANGER":      3,
		"OFFLINE":     4,
		"MAINTENANCE": 5,
	}
)

func (x ServiceState_Status) Enum() *ServiceState_Status {
	p := new(ServiceState_Status)
	*p = x
	return p
}

func (x ServiceState_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServiceState_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_trisa_api_v1beta1_api_proto_enumTypes[0].Descriptor()
}

func (ServiceState_Status) Type() protoreflect.EnumType {
	return &file_trisa_api_v1beta1_api_proto_enumTypes[0]
}

func (x ServiceState_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServiceState_Status.Descriptor instead.
func (ServiceState_Status) EnumDescriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{6, 0}
}

// Encrypted transaction envelope that is the outer layer of the TRISA information
// exchange protocol and facilitates the secure storage of KYC data in a transaction.
// The envelope specifies a unique id to reference the transaction out-of-band (e.g in
// the blockchain layer) and provides the necessary information so that only the
// originator and the beneficiary can decrypt the trnasaction data.
type SecureEnvelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The transaction identifier generated by the sender. Any response
	// to a transaction request needs to carry the same identifier.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Encrypted Payload
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	// Encryption key used to encrypt the transaction blob. This key itself
	// is encrypted using the public key of the receiver.
	EncryptionKey []byte `protobuf:"bytes,3,opt,name=encryption_key,json=encryptionKey,proto3" json:"encryption_key,omitempty"`
	// The encryption algorithm used to encrypt the transaction blob.
	EncryptionAlgorithm string `protobuf:"bytes,4,opt,name=encryption_algorithm,json=encryptionAlgorithm,proto3" json:"encryption_algorithm,omitempty"`
	// HMAC signature calculated from encrypted transaction blob.
	Hmac []byte `protobuf:"bytes,5,opt,name=hmac,proto3" json:"hmac,omitempty"`
	// The HMAC secret used to calculate the HMAC signature. This secret
	// itself is encrypted using the public key of the receiver.
	HmacSecret []byte `protobuf:"bytes,6,opt,name=hmac_secret,json=hmacSecret,proto3" json:"hmac_secret,omitempty"`
	// The algorithm used to calculate the HMAC signature.
	HmacAlgorithm string `protobuf:"bytes,7,opt,name=hmac_algorithm,json=hmacAlgorithm,proto3" json:"hmac_algorithm,omitempty"`
	// Rejection errors are used inside of a streaming context so that the stream is
	// not closed when an exchange-related rejection occurs. In the unary case, errors
	// are directly returned rather than as part of the secure envelope.
	Error *Error `protobuf:"bytes,9,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SecureEnvelope) Reset() {
	*x = SecureEnvelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecureEnvelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecureEnvelope) ProtoMessage() {}

func (x *SecureEnvelope) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecureEnvelope.ProtoReflect.Descriptor instead.
func (*SecureEnvelope) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{0}
}

func (x *SecureEnvelope) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SecureEnvelope) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *SecureEnvelope) GetEncryptionKey() []byte {
	if x != nil {
		return x.EncryptionKey
	}
	return nil
}

func (x *SecureEnvelope) GetEncryptionAlgorithm() string {
	if x != nil {
		return x.EncryptionAlgorithm
	}
	return ""
}

func (x *SecureEnvelope) GetHmac() []byte {
	if x != nil {
		return x.Hmac
	}
	return nil
}

func (x *SecureEnvelope) GetHmacSecret() []byte {
	if x != nil {
		return x.HmacSecret
	}
	return nil
}

func (x *SecureEnvelope) GetHmacAlgorithm() string {
	if x != nil {
		return x.HmacAlgorithm
	}
	return ""
}

func (x *SecureEnvelope) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

// Payload is the expected message structure that will be parsed from the encrypted
// secure envelope. The Payload should contain the identity and transaction information
// for the information exchange. The internal message types are purposefully generic to
// allow flexibility with the data needs for different exchanges.
type Payload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Identity contains any valid identity structure. The recommended format is the
	// IVMS101 IdentityPayload which contains the originator and beneficiary identities,
	// the originator and beneficiary VASP identities, as well as the transfer path of
	// any intermediate VASPs. The identity payload can be bidirectional (containing
	// both originator and beneficiary identities) or unidirectional - containing only
	// the identity of the sender. In the bidirectional case, the identity may be
	// purposefully partial to allow the recipient to fill in the details. In the
	// unidirectional case, the identities must be collated after.
	Identity *anypb.Any `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	// Transaction contains network specific information about the exchange or transfer.
	Transaction *anypb.Any `protobuf:"bytes,2,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *Payload) Reset() {
	*x = Payload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payload) ProtoMessage() {}

func (x *Payload) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payload.ProtoReflect.Descriptor instead.
func (*Payload) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{1}
}

func (x *Payload) GetIdentity() *anypb.Any {
	if x != nil {
		return x.Identity
	}
	return nil
}

func (x *Payload) GetTransaction() *anypb.Any {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// TODO: specify the address confirmation protocol.
type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{2}
}

type AddressConfirmation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddressConfirmation) Reset() {
	*x = AddressConfirmation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddressConfirmation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressConfirmation) ProtoMessage() {}

func (x *AddressConfirmation) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressConfirmation.ProtoReflect.Descriptor instead.
func (*AddressConfirmation) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{3}
}

// SigningKey provides metadata for decoding a PEM encoded PKIX public key for RSA
// encryption and transaction signing. The SigningKey is a lightweight version of the
// Certificate information stored in the Directory Service.
type SigningKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// x.509 metadata information for ease of reference without parsing the key.
	Version            int64  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Signature          []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	SignatureAlgorithm string `protobuf:"bytes,3,opt,name=signature_algorithm,json=signatureAlgorithm,proto3" json:"signature_algorithm,omitempty"`
	PublicKeyAlgorithm string `protobuf:"bytes,4,opt,name=public_key_algorithm,json=publicKeyAlgorithm,proto3" json:"public_key_algorithm,omitempty"`
	// Validity information
	NotBefore string `protobuf:"bytes,8,opt,name=not_before,json=notBefore,proto3" json:"not_before,omitempty"`
	NotAfter  string `protobuf:"bytes,9,opt,name=not_after,json=notAfter,proto3" json:"not_after,omitempty"`
	Revoked   bool   `protobuf:"varint,10,opt,name=revoked,proto3" json:"revoked,omitempty"`
	// The PEM encoded public key to PKIX, ASN.1 DER form without the trust chain.
	Data []byte `protobuf:"bytes,11,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *SigningKey) Reset() {
	*x = SigningKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SigningKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SigningKey) ProtoMessage() {}

func (x *SigningKey) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SigningKey.ProtoReflect.Descriptor instead.
func (*SigningKey) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{4}
}

func (x *SigningKey) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *SigningKey) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *SigningKey) GetSignatureAlgorithm() string {
	if x != nil {
		return x.SignatureAlgorithm
	}
	return ""
}

func (x *SigningKey) GetPublicKeyAlgorithm() string {
	if x != nil {
		return x.PublicKeyAlgorithm
	}
	return ""
}

func (x *SigningKey) GetNotBefore() string {
	if x != nil {
		return x.NotBefore
	}
	return ""
}

func (x *SigningKey) GetNotAfter() string {
	if x != nil {
		return x.NotAfter
	}
	return ""
}

func (x *SigningKey) GetRevoked() bool {
	if x != nil {
		return x.Revoked
	}
	return false
}

func (x *SigningKey) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type HealthCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The number of failed health checks that proceeded the current check.
	Attempts uint32 `protobuf:"varint,1,opt,name=attempts,proto3" json:"attempts,omitempty"`
	// The timestamp of the last health check, successful or otherwise.
	LastCheckedAt string `protobuf:"bytes,2,opt,name=last_checked_at,json=lastCheckedAt,proto3" json:"last_checked_at,omitempty"`
}

func (x *HealthCheck) Reset() {
	*x = HealthCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheck) ProtoMessage() {}

func (x *HealthCheck) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheck.ProtoReflect.Descriptor instead.
func (*HealthCheck) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{5}
}

func (x *HealthCheck) GetAttempts() uint32 {
	if x != nil {
		return x.Attempts
	}
	return 0
}

func (x *HealthCheck) GetLastCheckedAt() string {
	if x != nil {
		return x.LastCheckedAt
	}
	return ""
}

type ServiceState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Current service status as defined by the recieving system. The system is obliged
	// to respond with the closest matching status in a best-effort fashion. Alerts will
	// be triggered on service status changes if the system does not respond and the
	// previous system state was not unknown.
	Status ServiceState_Status `protobuf:"varint,1,opt,name=status,proto3,enum=trisa.api.v1beta1.ServiceState_Status" json:"status,omitempty"`
	// Suggest to the directory service when to check the health status again.
	NotBefore string `protobuf:"bytes,2,opt,name=not_before,json=notBefore,proto3" json:"not_before,omitempty"`
	NotAfter  string `protobuf:"bytes,3,opt,name=not_after,json=notAfter,proto3" json:"not_after,omitempty"`
}

func (x *ServiceState) Reset() {
	*x = ServiceState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_api_v1beta1_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceState) ProtoMessage() {}

func (x *ServiceState) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_api_v1beta1_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceState.ProtoReflect.Descriptor instead.
func (*ServiceState) Descriptor() ([]byte, []int) {
	return file_trisa_api_v1beta1_api_proto_rawDescGZIP(), []int{6}
}

func (x *ServiceState) GetStatus() ServiceState_Status {
	if x != nil {
		return x.Status
	}
	return ServiceState_UNKNOWN
}

func (x *ServiceState) GetNotBefore() string {
	if x != nil {
		return x.NotBefore
	}
	return ""
}

func (x *ServiceState) GetNotAfter() string {
	if x != nil {
		return x.NotAfter
	}
	return ""
}

var File_trisa_api_v1beta1_api_proto protoreflect.FileDescriptor

var file_trisa_api_v1beta1_api_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x74,
	0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x74, 0x72, 0x69,
	0x73, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x02, 0x0a, 0x0e,
	0x53, 0x65, 0x63, 0x75, 0x72, 0x65, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0d, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x12,
	0x31, 0x0a, 0x14, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x61, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74,
	0x68, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6d, 0x61, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x68, 0x6d, 0x61, 0x63, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x6d, 0x61, 0x63, 0x5f, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x68, 0x6d, 0x61,
	0x63, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x68, 0x6d, 0x61, 0x63, 0x5f,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x68, 0x6d, 0x61, 0x63, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x2e,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x73,
	0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x30, 0x0a, 0x08, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x36, 0x0a, 0x0b, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x09, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x15,
	0x0a, 0x13, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x91, 0x02, 0x0a, 0x0a, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e,
	0x67, 0x4b, 0x65, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x2f, 0x0a, 0x13,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69,
	0x74, 0x68, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x30, 0x0a,
	0x14, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x61, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12,
	0x1d, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x6f, 0x74, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x6e, 0x6f, 0x74, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x6f, 0x74, 0x41, 0x66, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x65, 0x76, 0x6f, 0x6b, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65,
	0x76, 0x6f, 0x6b, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x51, 0x0a, 0x0b, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x74, 0x74, 0x65,
	0x6d, 0x70, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x61, 0x74, 0x74, 0x65,
	0x6d, 0x70, 0x74, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x63, 0x68, 0x65,
	0x63, 0x6b, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c,
	0x61, 0x73, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x65, 0x64, 0x41, 0x74, 0x22, 0xe7, 0x01, 0x0a,
	0x0c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3e, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e,
	0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a,
	0x0a, 0x6e, 0x6f, 0x74, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6e, 0x6f, 0x74, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x6e, 0x6f, 0x74, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6e, 0x6f, 0x74, 0x41, 0x66, 0x74, 0x65, 0x72, 0x22, 0x5b, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x01, 0x12, 0x0d, 0x0a,
	0x09, 0x55, 0x4e, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06,
	0x44, 0x41, 0x4e, 0x47, 0x45, 0x52, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x4f, 0x46, 0x46, 0x4c,
	0x49, 0x4e, 0x45, 0x10, 0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x41, 0x49, 0x4e, 0x54, 0x45, 0x4e,
	0x41, 0x4e, 0x43, 0x45, 0x10, 0x05, 0x32, 0xe7, 0x02, 0x0a, 0x0c, 0x54, 0x52, 0x49, 0x53, 0x41,
	0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x52, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x65, 0x72, 0x12, 0x21, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x65, 0x45, 0x6e,
	0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x1a, 0x21, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72,
	0x65, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x22, 0x00, 0x12, 0x5c, 0x0a, 0x0e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x21, 0x2e,
	0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x65, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65,
	0x1a, 0x21, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x65, 0x45, 0x6e, 0x76, 0x65, 0x6c,
	0x6f, 0x70, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x56, 0x0a, 0x0e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x72, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x2e, 0x74, 0x72,
	0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x1a, 0x26, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x00, 0x12, 0x4d, 0x0a, 0x0b, 0x4b, 0x65, 0x79, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x1d, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x1a,
	0x1d, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x22, 0x00,
	0x32, 0x5a, 0x0a, 0x0b, 0x54, 0x52, 0x49, 0x53, 0x41, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12,
	0x4b, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x2e, 0x74, 0x72, 0x69, 0x73,
	0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x48, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x1a, 0x1f, 0x2e, 0x74, 0x72, 0x69, 0x73,
	0x61, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0x00, 0x42, 0x38, 0x5a, 0x36,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x69, 0x73, 0x61,
	0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trisa_api_v1beta1_api_proto_rawDescOnce sync.Once
	file_trisa_api_v1beta1_api_proto_rawDescData = file_trisa_api_v1beta1_api_proto_rawDesc
)

func file_trisa_api_v1beta1_api_proto_rawDescGZIP() []byte {
	file_trisa_api_v1beta1_api_proto_rawDescOnce.Do(func() {
		file_trisa_api_v1beta1_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_trisa_api_v1beta1_api_proto_rawDescData)
	})
	return file_trisa_api_v1beta1_api_proto_rawDescData
}

var file_trisa_api_v1beta1_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_trisa_api_v1beta1_api_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_trisa_api_v1beta1_api_proto_goTypes = []interface{}{
	(ServiceState_Status)(0),    // 0: trisa.api.v1beta1.ServiceState.Status
	(*SecureEnvelope)(nil),      // 1: trisa.api.v1beta1.SecureEnvelope
	(*Payload)(nil),             // 2: trisa.api.v1beta1.Payload
	(*Address)(nil),             // 3: trisa.api.v1beta1.Address
	(*AddressConfirmation)(nil), // 4: trisa.api.v1beta1.AddressConfirmation
	(*SigningKey)(nil),          // 5: trisa.api.v1beta1.SigningKey
	(*HealthCheck)(nil),         // 6: trisa.api.v1beta1.HealthCheck
	(*ServiceState)(nil),        // 7: trisa.api.v1beta1.ServiceState
	(*Error)(nil),               // 8: trisa.api.v1beta1.Error
	(*anypb.Any)(nil),           // 9: google.protobuf.Any
}
var file_trisa_api_v1beta1_api_proto_depIdxs = []int32{
	8, // 0: trisa.api.v1beta1.SecureEnvelope.error:type_name -> trisa.api.v1beta1.Error
	9, // 1: trisa.api.v1beta1.Payload.identity:type_name -> google.protobuf.Any
	9, // 2: trisa.api.v1beta1.Payload.transaction:type_name -> google.protobuf.Any
	0, // 3: trisa.api.v1beta1.ServiceState.status:type_name -> trisa.api.v1beta1.ServiceState.Status
	1, // 4: trisa.api.v1beta1.TRISANetwork.Transfer:input_type -> trisa.api.v1beta1.SecureEnvelope
	1, // 5: trisa.api.v1beta1.TRISANetwork.TransferStream:input_type -> trisa.api.v1beta1.SecureEnvelope
	3, // 6: trisa.api.v1beta1.TRISANetwork.ConfirmAddress:input_type -> trisa.api.v1beta1.Address
	5, // 7: trisa.api.v1beta1.TRISANetwork.KeyExchange:input_type -> trisa.api.v1beta1.SigningKey
	6, // 8: trisa.api.v1beta1.TRISAHealth.Status:input_type -> trisa.api.v1beta1.HealthCheck
	1, // 9: trisa.api.v1beta1.TRISANetwork.Transfer:output_type -> trisa.api.v1beta1.SecureEnvelope
	1, // 10: trisa.api.v1beta1.TRISANetwork.TransferStream:output_type -> trisa.api.v1beta1.SecureEnvelope
	4, // 11: trisa.api.v1beta1.TRISANetwork.ConfirmAddress:output_type -> trisa.api.v1beta1.AddressConfirmation
	5, // 12: trisa.api.v1beta1.TRISANetwork.KeyExchange:output_type -> trisa.api.v1beta1.SigningKey
	7, // 13: trisa.api.v1beta1.TRISAHealth.Status:output_type -> trisa.api.v1beta1.ServiceState
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_trisa_api_v1beta1_api_proto_init() }
func file_trisa_api_v1beta1_api_proto_init() {
	if File_trisa_api_v1beta1_api_proto != nil {
		return
	}
	file_trisa_api_v1beta1_errors_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_trisa_api_v1beta1_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecureEnvelope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trisa_api_v1beta1_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Payload); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trisa_api_v1beta1_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trisa_api_v1beta1_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddressConfirmation); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trisa_api_v1beta1_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SigningKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trisa_api_v1beta1_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheck); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_trisa_api_v1beta1_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_trisa_api_v1beta1_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_trisa_api_v1beta1_api_proto_goTypes,
		DependencyIndexes: file_trisa_api_v1beta1_api_proto_depIdxs,
		EnumInfos:         file_trisa_api_v1beta1_api_proto_enumTypes,
		MessageInfos:      file_trisa_api_v1beta1_api_proto_msgTypes,
	}.Build()
	File_trisa_api_v1beta1_api_proto = out.File
	file_trisa_api_v1beta1_api_proto_rawDesc = nil
	file_trisa_api_v1beta1_api_proto_goTypes = nil
	file_trisa_api_v1beta1_api_proto_depIdxs = nil
}
