package api

import (
	"fmt"

	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// TRISA error code constants. See protocol buffers documentation for more details.
const (
	Unhandled                = api.Error_UNHANDLED
	Unavailable              = api.Error_UNAVAILABLE
	ServiceDownTime          = api.Error_SERVICE_DOWN_TIME
	Maintenance              = api.Error_MAINTENANCE
	Unimplemented            = api.Error_UNIMPLEMENTED
	InternalError            = api.Error_INTERNAL_ERROR
	Rejected                 = api.Error_REJECTED
	UnkownWalletAddress      = api.Error_UNKNOWN_WALLET_ADDRESS
	UnknownIdentity          = api.Error_UNKOWN_IDENTITY
	UnkownOriginator         = api.Error_UNKNOWN_ORIGINATOR
	UnkownBeneficiary        = api.Error_UNKOWN_BENEFICIARY
	BeneficiaryNameUnmatched = api.Error_BENEFICIARY_NAME_UNMATCHED
	UnsupportedCurrency      = api.Error_UNSUPPORTED_CURRENCY
	ExceededTradingVolume    = api.Error_EXCEEDED_TRADING_VOLUME
	ComplianceCheckFail      = api.Error_COMPLIANCE_CHECK_FAIL
	NoCompliance             = api.Error_NO_COMPLIANCE
	HighRisk                 = api.Error_HIGH_RISK
	OutOfNetwork             = api.Error_OUT_OF_NETWORK
	Forbidden                = api.Error_FORBIDDEN
	NoSigningKey             = api.Error_NO_SIGNING_KEY
	CertificateRevoked       = api.Error_CERTIFICATE_REVOKED
	Unverified               = api.Error_UNVERIFIED
	Untrusted                = api.Error_UNTRUSTED
	InvalidSignature         = api.Error_INVALID_SIGNATURE
	InvalidKey               = api.Error_INVALID_KEY
	EnvelopeDecodeFail       = api.Error_ENVELOPE_DECODE_FAIL
	PrivateInfoDecodeFail    = api.Error_PRIVATE_INFO_DECODE_FAIL
	UnhandledAlgorithm       = api.Error_UNHANDLED_ALGORITHM
	BadRequest               = api.Error_BAD_REQUEST
	UnparseableIdentity      = api.Error_UNPARSEABLE_IDENTITY
	PrivateInfoWrongFormat   = api.Error_PRIVATE_INFO_WRONG_FORMAT
	UnparseableTransaction   = api.Error_UNPARSEABLE_TRANSACTION
	MissingFields            = api.Error_MISSING_FIELDS
	IncompleteIdentity       = api.Error_INCOMPLETE_IDENTITY
	ValidationError          = api.Error_VALIDATION_ERROR
)

// Sygna BVRC rejected error codes
const (
	BVRC001 = api.Error_BVRC001
	BVRC002 = api.Error_BVRC002
	BVRC003 = api.Error_BVRC003
	BVRC004 = api.Error_BVRC004
	BVRC005 = api.Error_BVRC005
	BVRC006 = api.Error_BVRC006
	BVRC007 = api.Error_BVRC007
	BVRC999 = api.Error_BVRC999
)

// Errorp parses an error from a status error (e.g. the error is embedded in the details)
// or if the error is already the correct type return that directly. If the error cannot
// be parsed, an *api.Error is returned with the Unhandled error code and the message of the
// original error; in this case, ok will be false.
func Errorp(err error) (e *api.Error, ok bool) {
	if err == nil {
		return nil, true
	}

	var s *status.Status
	if s, ok = status.FromError(err); ok {
		es := s.Details()
		if len(es) == 1 {
			if e, ok = es[0].(*api.Error); ok {
				return e, ok
			}
		}
	}

	return &api.Error{Code: Unhandled, Message: err.Error()}, false
}

// Errorf creates a new Error message formated with the specified arguments. If the
// error code indicates the error should be retried it sets retry to true.
func Errorf(code api.Error_Code, format string, a ...interface{}) *api.Error {
	if len(a) == 0 {
		return &api.Error{Code: code, Message: format}
	}

	return &api.Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Retry:   code >= 150,
	}
}

// WithRetry returns a new error with the retry flag set to true.
func WithRetry(e *api.Error) *api.Error {
	return &api.Error{
		Code:    e.Code,
		Message: e.Message,
		Retry:   true,
	}
}

// WithDetails returns the new error with the details as a pb.Any struct.
func WithDetails(e *api.Error, details proto.Message) (_ *api.Error, err error) {
	var any *anypb.Any
	if any, err = anypb.New(details); err != nil {
		return nil, err
	}
	return &api.Error{
		Code:    e.Code,
		Message: e.Message,
		Retry:   e.Retry,
		Details: any,
	}, nil
}

// IsZero returns true if the error has a code == 0 and no message
func IsZero(e *api.Error) bool {
	return e.Code == 0 && e.Message == ""
}

// Err returns a gRPC status error with appropriate gRPC status codes for returning
// out of a gRPC server function. This should be returned where possible.
// Deprecated: return an error inside of the envelope for a rejection.
func Err(e *api.Error) (err error) {
	if e == nil {
		return nil
	}

	var code codes.Code
	switch {
	case e.Code < 49:
		code = codes.Unavailable
	case e.Code == 49:
		code = codes.Internal
	case e.Code > 49 && e.Code < 100:
		code = codes.Aborted
	case e.Code >= 100 && e.Code < 105:
		code = codes.FailedPrecondition
	case e.Code > 105:
		code = codes.InvalidArgument
	}

	st := status.New(code, fmt.Sprintf("[%s] %s", e.Code, e.Message))
	if st, err = st.WithDetails(e); err != nil {
		panic(fmt.Sprintf("error attaching status metadata: %v", err))
	}
	return st.Err()
}
