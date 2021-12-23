package api

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// TRISA error code constants. See protocol buffers documentation for more details.
const (
	Unhandled                = Error_UNHANDLED
	Unavailable              = Error_UNAVAILABLE
	ServiceDownTime          = Error_SERVICE_DOWN_TIME
	Maintenance              = Error_MAINTENANCE
	Unimplemented            = Error_UNIMPLEMENTED
	InternalError            = Error_INTERNAL_ERROR
	Rejected                 = Error_REJECTED
	UnkownWalletAddress      = Error_UNKNOWN_WALLET_ADDRESS
	UnknownIdentity          = Error_UNKOWN_IDENTITY
	UnkownOriginator         = Error_UNKNOWN_ORIGINATOR
	UnkownBeneficiary        = Error_UNKOWN_BENEFICIARY
	BeneficiaryNameUnmatched = Error_BENEFICIARY_NAME_UNMATCHED
	UnsupportedCurrency      = Error_UNSUPPORTED_CURRENCY
	ExceededTradingVolume    = Error_EXCEEDED_TRADING_VOLUME
	ComplianceCheckFail      = Error_COMPLIANCE_CHECK_FAIL
	NoCompliance             = Error_NO_COMPLIANCE
	HighRisk                 = Error_HIGH_RISK
	OutOfNetwork             = Error_OUT_OF_NETWORK
	Forbidden                = Error_FORBIDDEN
	NoSigningKey             = Error_NO_SIGNING_KEY
	CertificateRevoked       = Error_CERTIFICATE_REVOKED
	Unverified               = Error_UNVERIFIED
	Untrusted                = Error_UNTRUSTED
	InvalidSignature         = Error_INVALID_SIGNATURE
	InvalidKey               = Error_INVALID_KEY
	EnvelopeDecodeFail       = Error_ENVELOPE_DECODE_FAIL
	PrivateInfoDecodeFail    = Error_PRIVATE_INFO_DECODE_FAIL
	UnhandledAlgorithm       = Error_UNHANDLED_ALGORITHM
	BadRequest               = Error_BAD_REQUEST
	UnparseableIdentity      = Error_UNPARSEABLE_IDENTITY
	PrivateInfoWrongFormat   = Error_PRIVATE_INFO_WRONG_FORMAT
	UnparseableTransaction   = Error_UNPARSEABLE_TRANSACTION
	MissingFields            = Error_MISSING_FIELDS
	IncompleteIdentity       = Error_INCOMPLETE_IDENTITY
	ValidationError          = Error_VALIDATION_ERROR
)

// Sygna BVRC rejected error codes
const (
	BVRC001 = Error_BVRC001
	BVRC002 = Error_BVRC002
	BVRC003 = Error_BVRC003
	BVRC004 = Error_BVRC004
	BVRC005 = Error_BVRC005
	BVRC006 = Error_BVRC006
	BVRC007 = Error_BVRC007
	BVRC999 = Error_BVRC999
)

// Errorp parses an error from a status error (e.g. the error is embedded in the details)
// or if the error is already the correct type return that directly. If the error cannot
// be parsed, an *Error is returned with the Unhandled error code and the message of the
// original error; in this case, ok will be false.
func Errorp(err error) (e *Error, ok bool) {
	if err == nil {
		return nil, true
	}

	if e, ok = err.(*Error); ok {
		return e, ok
	}

	var s *status.Status
	if s, ok = status.FromError(err); ok {
		es := s.Details()
		if len(es) == 1 {
			if e, ok = es[0].(*Error); ok {
				return e, ok
			}
		}
	}

	return &Error{Code: Unhandled, Message: err.Error()}, false
}

// Errorf creates a new Error message formated with the specified arguments. If the
// error code indicates the error should be retried it sets retry to true.
func Errorf(code Error_Code, format string, a ...interface{}) *Error {
	if len(a) == 0 {
		return &Error{Code: code, Message: format}
	}

	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
		Retry:   code >= 150,
	}
}

// WithRetry returns a new error with the retry flag set to true.
func (e *Error) WithRetry() *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Retry:   true,
	}
}

// WithDetails returns the new error with the details as a pb.Any struct.
func (e *Error) WithDetails(details proto.Message) (_ *Error, err error) {
	var any *anypb.Any
	if any, err = anypb.New(details); err != nil {
		return nil, err
	}
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Retry:   e.Retry,
		Details: any,
	}, nil
}

// Error implements the error interface for printing and logging.
func (e *Error) Error() string {
	return fmt.Sprintf("trisa error [%s]: %s", e.Code, e.Message)
}

// Err returns a gRPC status error with appropriate gRPC status codes for returning
// out of a gRPC server function. This should be returned where possible.
func (e *Error) Err() (err error) {
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
