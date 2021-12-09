// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: trisa/gds/models/v1beta1/ca.proto

package models

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Certificate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// x.509 metadata information for ease of reference
	// note that the serial number as a hex string can be used with the Sectigo API
	Version            int64  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	SerialNumber       []byte `protobuf:"bytes,2,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	Signature          []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	SignatureAlgorithm string `protobuf:"bytes,4,opt,name=signature_algorithm,json=signatureAlgorithm,proto3" json:"signature_algorithm,omitempty"`
	PublicKeyAlgorithm string `protobuf:"bytes,5,opt,name=public_key_algorithm,json=publicKeyAlgorithm,proto3" json:"public_key_algorithm,omitempty"`
	// Issuer and subject information from Sectigo
	Subject *Name `protobuf:"bytes,6,opt,name=subject,proto3" json:"subject,omitempty"`
	Issuer  *Name `protobuf:"bytes,7,opt,name=issuer,proto3" json:"issuer,omitempty"`
	// Validity information
	NotBefore string `protobuf:"bytes,8,opt,name=not_before,json=notBefore,proto3" json:"not_before,omitempty"`
	NotAfter  string `protobuf:"bytes,9,opt,name=not_after,json=notAfter,proto3" json:"not_after,omitempty"`
	Revoked   bool   `protobuf:"varint,10,opt,name=revoked,proto3" json:"revoked,omitempty"`
	// The ASN1 encoded full certificate without the trust chain
	Data []byte `protobuf:"bytes,11,opt,name=data,proto3" json:"data,omitempty"`
	// The complete trust chain including the leaf certificate as a gzip compressed
	// PEM encoded file. This field can be deserialized into a trust.Provider.
	Chain []byte `protobuf:"bytes,12,opt,name=chain,proto3" json:"chain,omitempty"`
}

func (x *Certificate) Reset() {
	*x = Certificate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_gds_models_v1beta1_ca_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Certificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certificate) ProtoMessage() {}

func (x *Certificate) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_gds_models_v1beta1_ca_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certificate.ProtoReflect.Descriptor instead.
func (*Certificate) Descriptor() ([]byte, []int) {
	return file_trisa_gds_models_v1beta1_ca_proto_rawDescGZIP(), []int{0}
}

func (x *Certificate) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Certificate) GetSerialNumber() []byte {
	if x != nil {
		return x.SerialNumber
	}
	return nil
}

func (x *Certificate) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *Certificate) GetSignatureAlgorithm() string {
	if x != nil {
		return x.SignatureAlgorithm
	}
	return ""
}

func (x *Certificate) GetPublicKeyAlgorithm() string {
	if x != nil {
		return x.PublicKeyAlgorithm
	}
	return ""
}

func (x *Certificate) GetSubject() *Name {
	if x != nil {
		return x.Subject
	}
	return nil
}

func (x *Certificate) GetIssuer() *Name {
	if x != nil {
		return x.Issuer
	}
	return nil
}

func (x *Certificate) GetNotBefore() string {
	if x != nil {
		return x.NotBefore
	}
	return ""
}

func (x *Certificate) GetNotAfter() string {
	if x != nil {
		return x.NotAfter
	}
	return ""
}

func (x *Certificate) GetRevoked() bool {
	if x != nil {
		return x.Revoked
	}
	return false
}

func (x *Certificate) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Certificate) GetChain() []byte {
	if x != nil {
		return x.Chain
	}
	return nil
}

// An X.509 distinguished name with the common elements of a DN.
type Name struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonName         string   `protobuf:"bytes,1,opt,name=common_name,json=commonName,proto3" json:"common_name,omitempty"`
	SerialNumber       string   `protobuf:"bytes,2,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	Organization       []string `protobuf:"bytes,3,rep,name=organization,proto3" json:"organization,omitempty"`
	OrganizationalUnit []string `protobuf:"bytes,4,rep,name=organizational_unit,json=organizationalUnit,proto3" json:"organizational_unit,omitempty"`
	StreetAddress      []string `protobuf:"bytes,5,rep,name=street_address,json=streetAddress,proto3" json:"street_address,omitempty"`
	Locality           []string `protobuf:"bytes,6,rep,name=locality,proto3" json:"locality,omitempty"`
	Province           []string `protobuf:"bytes,7,rep,name=province,proto3" json:"province,omitempty"`
	PostalCode         []string `protobuf:"bytes,8,rep,name=postal_code,json=postalCode,proto3" json:"postal_code,omitempty"`
	Country            []string `protobuf:"bytes,9,rep,name=country,proto3" json:"country,omitempty"`
}

func (x *Name) Reset() {
	*x = Name{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trisa_gds_models_v1beta1_ca_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Name) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Name) ProtoMessage() {}

func (x *Name) ProtoReflect() protoreflect.Message {
	mi := &file_trisa_gds_models_v1beta1_ca_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Name.ProtoReflect.Descriptor instead.
func (*Name) Descriptor() ([]byte, []int) {
	return file_trisa_gds_models_v1beta1_ca_proto_rawDescGZIP(), []int{1}
}

func (x *Name) GetCommonName() string {
	if x != nil {
		return x.CommonName
	}
	return ""
}

func (x *Name) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *Name) GetOrganization() []string {
	if x != nil {
		return x.Organization
	}
	return nil
}

func (x *Name) GetOrganizationalUnit() []string {
	if x != nil {
		return x.OrganizationalUnit
	}
	return nil
}

func (x *Name) GetStreetAddress() []string {
	if x != nil {
		return x.StreetAddress
	}
	return nil
}

func (x *Name) GetLocality() []string {
	if x != nil {
		return x.Locality
	}
	return nil
}

func (x *Name) GetProvince() []string {
	if x != nil {
		return x.Province
	}
	return nil
}

func (x *Name) GetPostalCode() []string {
	if x != nil {
		return x.PostalCode
	}
	return nil
}

func (x *Name) GetCountry() []string {
	if x != nil {
		return x.Country
	}
	return nil
}

var File_trisa_gds_models_v1beta1_ca_proto protoreflect.FileDescriptor

var file_trisa_gds_models_v1beta1_ca_proto_rawDesc = []byte{
	0x0a, 0x21, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x18, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x22, 0xbf, 0x03,
	0x0a, 0x0b, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x2f, 0x0a, 0x13, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x30, 0x0a, 0x14, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69,
	0x74, 0x68, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x4b, 0x65, 0x79, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x38, 0x0a,
	0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e, 0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x07,
	0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2e,
	0x67, 0x64, 0x73, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x06, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x12,
	0x1d, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x6f, 0x74, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x6e, 0x6f, 0x74, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6e, 0x6f, 0x74, 0x41, 0x66, 0x74, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x65, 0x76, 0x6f, 0x6b, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65,
	0x76, 0x6f, 0x6b, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x68, 0x61,
	0x69, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x22,
	0xbb, 0x02, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72,
	0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x22,
	0x0a, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x13, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x12, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x55,
	0x6e, 0x69, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74, 0x72, 0x65, 0x65, 0x74, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x72,
	0x65, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x6e,
	0x63, 0x65, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x6e,
	0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6f, 0x73, 0x74, 0x61, 0x6c, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x09,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x42, 0x5a,
	0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x69, 0x73,
	0x61, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x74, 0x72, 0x69, 0x73, 0x61, 0x2f, 0x67, 0x64, 0x73, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trisa_gds_models_v1beta1_ca_proto_rawDescOnce sync.Once
	file_trisa_gds_models_v1beta1_ca_proto_rawDescData = file_trisa_gds_models_v1beta1_ca_proto_rawDesc
)

func file_trisa_gds_models_v1beta1_ca_proto_rawDescGZIP() []byte {
	file_trisa_gds_models_v1beta1_ca_proto_rawDescOnce.Do(func() {
		file_trisa_gds_models_v1beta1_ca_proto_rawDescData = protoimpl.X.CompressGZIP(file_trisa_gds_models_v1beta1_ca_proto_rawDescData)
	})
	return file_trisa_gds_models_v1beta1_ca_proto_rawDescData
}

var file_trisa_gds_models_v1beta1_ca_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_trisa_gds_models_v1beta1_ca_proto_goTypes = []interface{}{
	(*Certificate)(nil), // 0: trisa.gds.models.v1beta1.Certificate
	(*Name)(nil),        // 1: trisa.gds.models.v1beta1.Name
}
var file_trisa_gds_models_v1beta1_ca_proto_depIdxs = []int32{
	1, // 0: trisa.gds.models.v1beta1.Certificate.subject:type_name -> trisa.gds.models.v1beta1.Name
	1, // 1: trisa.gds.models.v1beta1.Certificate.issuer:type_name -> trisa.gds.models.v1beta1.Name
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_trisa_gds_models_v1beta1_ca_proto_init() }
func file_trisa_gds_models_v1beta1_ca_proto_init() {
	if File_trisa_gds_models_v1beta1_ca_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trisa_gds_models_v1beta1_ca_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Certificate); i {
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
		file_trisa_gds_models_v1beta1_ca_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Name); i {
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
			RawDescriptor: file_trisa_gds_models_v1beta1_ca_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trisa_gds_models_v1beta1_ca_proto_goTypes,
		DependencyIndexes: file_trisa_gds_models_v1beta1_ca_proto_depIdxs,
		MessageInfos:      file_trisa_gds_models_v1beta1_ca_proto_msgTypes,
	}.Build()
	File_trisa_gds_models_v1beta1_ca_proto = out.File
	file_trisa_gds_models_v1beta1_ca_proto_rawDesc = nil
	file_trisa_gds_models_v1beta1_ca_proto_goTypes = nil
	file_trisa_gds_models_v1beta1_ca_proto_depIdxs = nil
}
