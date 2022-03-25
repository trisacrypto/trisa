package api_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestErrors(t *testing.T) {
	err := api.Errorf(api.UnknownIdentity, "could not parse %q", "foo")
	require.Error(t, err)
	require.Equal(t, err.Error(), `trisa rejection [UNKOWN_IDENTITY]: could not parse "foo"`)
	require.False(t, err.IsZero())

	oerr, ok := api.Errorp(err)
	require.True(t, ok)
	require.Equal(t, err, oerr)
	require.False(t, oerr.IsZero())

	oerr, ok = api.Errorp(errors.New("unhandled error"))
	require.False(t, ok)
	require.Equal(t, oerr.Error(), "trisa rejection [UNHANDLED]: unhandled error")
	require.False(t, oerr.IsZero())

	sterr := err.Err()
	require.Equal(t, sterr.Error(), `rpc error: code = Aborted desc = [UNKOWN_IDENTITY] could not parse "foo"`)

	oerr, ok = api.Errorp(sterr)
	require.True(t, ok)
	require.True(t, proto.Equal(err, oerr), "unexpected return value from Errorp")
	require.False(t, oerr.IsZero())

	// WithRetry should return a new error with retry set to true
	errWithRetry := err.WithRetry()
	require.Equal(t, err.Code, errWithRetry.Code)
	require.Equal(t, err.Message, errWithRetry.Message)
	require.True(t, errWithRetry.Retry)
	require.Nil(t, errWithRetry.Details)
	require.False(t, errWithRetry.IsZero())

	_, parseErr := err.WithDetails(nil)
	require.Error(t, parseErr)

	// WithDetails should add an arbitrary proto.Message to the error details
	details := &api.Error{
		Code: api.UnknownIdentity,
	}
	errWithDetails, parseErr := err.WithDetails(details)
	require.NoError(t, parseErr)
	require.Equal(t, err.Code, errWithDetails.Code)
	require.Equal(t, err.Message, errWithDetails.Message)
	require.Equal(t, err.Retry, errWithDetails.Retry)
	actualDetails := &api.Error{}
	require.NoError(t, anypb.UnmarshalTo(errWithDetails.Details, actualDetails, proto.UnmarshalOptions{}))
	require.True(t, proto.Equal(details, actualDetails), "unexpected details created by WithDetails")
	require.False(t, errWithDetails.IsZero())
}

func TestIsZero(t *testing.T) {
	err := &api.Error{}
	require.True(t, err.IsZero(), "no code and no message should be zero valued")

	err = &api.Error{Retry: true}
	require.True(t, err.IsZero(), "non-zero retry is not sufficient")

	details, _ := anypb.New(&api.Error{Code: api.ExceededTradingVolume, Message: "too fast"})
	err = &api.Error{Details: details}
	require.True(t, err.IsZero(), "non-zero details is not sufficient")

	err = &api.Error{Code: api.OutOfNetwork}
	require.False(t, err.IsZero(), "a code greater than zero should be sufficient")

	err = &api.Error{Message: "unexpected content"}
	require.False(t, err.IsZero(), "a message without a code should be sufficient")

	err = &api.Error{Code: api.OutOfNetwork, Message: "unexpected content"}
	require.False(t, err.IsZero(), "both a message and a code should be non-zero")

	// After marshaling an empty protocol buffer message, it should still be zero
	data, merr := proto.Marshal(&api.Error{})
	require.NoError(t, merr, "could not marshal protocol buffer")
	umerr := &api.Error{}
	require.NoError(t, proto.Unmarshal(data, umerr), "could not unmarshal error pb")
	require.True(t, umerr.IsZero(), "should be zero after marshal and unmarshal")
}

// Test that the Err function returns an error that includes the corresponding gRPC
// error code and message from the original error.
func TestErr(t *testing.T) {
	tests := []struct {
		name         string
		code         api.Error_Code
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
			err := &api.Error{
				Code:    tt.code,
				Message: tt.message,
			}
			stErr := err.Err()
			require.NotNil(t, stErr)
			require.Equal(t, tt.expectedCode, status.Code(stErr))
			require.Contains(t, stErr.Error(), tt.message)
		})
	}
}
