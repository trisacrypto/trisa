package api_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	trisa "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1/errors"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestErrors(t *testing.T) {
	err := api.Errorf(api.UnknownIdentity, "could not parse %q", "foo")
	require.False(t, api.IsZero(err))

	oerr, ok := api.Errorp(errors.New("unhandled error"))
	require.False(t, ok)
	require.False(t, api.IsZero(oerr))

	sterr := api.Err(err)
	require.Equal(t, sterr.Error(), `rpc error: code = Aborted desc = [UNKNOWN_IDENTITY] could not parse "foo"`)

	oerr, ok = api.Errorp(sterr)
	require.True(t, ok)
	require.True(t, proto.Equal(err, oerr), "unexpected return value from Errorp")
	require.False(t, api.IsZero(oerr))

	// WithRetry should return a new error with retry set to true
	errWithRetry := api.WithRetry(err)
	require.Equal(t, err.Code, errWithRetry.Code)
	require.Equal(t, err.Message, errWithRetry.Message)
	require.True(t, errWithRetry.Retry)
	require.Nil(t, errWithRetry.Details)
	require.False(t, api.IsZero(errWithRetry))

	_, parseErr := api.WithDetails(err, nil)
	require.Error(t, parseErr)

	// WithDetails should add an arbitrary proto.Message to the error details
	details := &trisa.Error{
		Code: api.UnknownIdentity,
	}
	errWithDetails, parseErr := api.WithDetails(err, details)
	require.NoError(t, parseErr)
	require.Equal(t, err.Code, errWithDetails.Code)
	require.Equal(t, err.Message, errWithDetails.Message)
	require.Equal(t, err.Retry, errWithDetails.Retry)
	actualDetails := &trisa.Error{}
	require.NoError(t, anypb.UnmarshalTo(errWithDetails.Details, actualDetails, proto.UnmarshalOptions{}))
	require.True(t, proto.Equal(details, actualDetails), "unexpected details created by WithDetails")
	require.False(t, api.IsZero(errWithDetails))
}

func TestIsZero(t *testing.T) {
	err := &trisa.Error{}
	require.True(t, api.IsZero(err), "no code and no message should be zero valued")

	err = &trisa.Error{Retry: true}
	require.True(t, api.IsZero(err), "non-zero retry is not sufficient")

	details, _ := anypb.New(&trisa.Error{Code: api.ExceededTradingVolume, Message: "too fast"})
	err = &trisa.Error{Details: details}
	require.True(t, api.IsZero(err), "non-zero details is not sufficient")

	err = &trisa.Error{Code: api.OutOfNetwork}
	require.False(t, api.IsZero(err), "a code greater than zero should be sufficient")

	err = &trisa.Error{Message: "unexpected content"}
	require.False(t, api.IsZero(err), "a message without a code should be sufficient")

	err = &trisa.Error{Code: api.OutOfNetwork, Message: "unexpected content"}
	require.False(t, api.IsZero(err), "both a message and a code should be non-zero")

	// After marshaling an empty protocol buffer message, it should still be zero
	data, merr := proto.Marshal(&trisa.Error{})
	require.NoError(t, merr, "could not marshal protocol buffer")
	umerr := &trisa.Error{}
	require.NoError(t, proto.Unmarshal(data, umerr), "could not unmarshal error pb")
	require.True(t, api.IsZero(umerr), "should be zero after marshal and unmarshal")
}

// Test that the Err function returns an error that includes the corresponding gRPC
// error code and message from the original error.
func TestErr(t *testing.T) {
	tests := []struct {
		name         string
		code         trisa.Error_Code
		message      string
		expectedCode codes.Code
	}{
		{"unavailable", api.Unavailable, "endpoint is unavailable", codes.Unavailable},
		{"internal", api.InternalError, "internal error", codes.Internal},
		{"aborted", api.Rejected, "aborted", codes.Aborted},
		{"failedPrecondition", api.Unverified, "failed precondition", codes.FailedPrecondition},
		{"invalidArgument", api.InvalidKey, "invalid argument", codes.InvalidArgument},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &trisa.Error{
				Code:    tt.code,
				Message: tt.message,
			}
			stErr := api.Err(err)
			require.NotNil(t, stErr)
			require.Equal(t, tt.expectedCode, status.Code(stErr))
			require.Contains(t, stErr.Error(), tt.message)
		})
	}
}
