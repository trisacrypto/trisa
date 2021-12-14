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
	require.Equal(t, err.Error(), `trisa error [UNKOWN_IDENTITY]: could not parse "foo"`)

	oerr, ok := api.Errorp(err)
	require.True(t, ok)
	require.Equal(t, err, oerr)

	oerr, ok = api.Errorp(errors.New("unhandled error"))
	require.False(t, ok)
	require.Equal(t, oerr.Error(), "trisa error [UNHANDLED]: unhandled error")

	sterr := err.Err()
	require.Equal(t, sterr.Error(), `rpc error: code = Aborted desc = [UNKOWN_IDENTITY] could not parse "foo"`)

	oerr, ok = api.Errorp(sterr)
	require.True(t, ok)
	require.True(t, proto.Equal(err, oerr), "unexpected return value from Errorp")

	// WithRetry should return a new error with retry set to true
	errWithRetry := err.WithRetry()
	require.Equal(t, err.Code, errWithRetry.Code)
	require.Equal(t, err.Message, errWithRetry.Message)
	require.True(t, errWithRetry.Retry)
	require.Nil(t, errWithRetry.Details)

	errWithDetails, parseErr := err.WithDetails(nil)
	require.Error(t, parseErr)

	// WithDetails should add an arbitrary proto.Message to the error details
	details := &api.Error{
		Code: api.UnknownIdentity,
	}
	errWithDetails, parseErr = err.WithDetails(details)
	require.Equal(t, err.Code, errWithDetails.Code)
	require.Equal(t, err.Message, errWithDetails.Message)
	require.Equal(t, err.Retry, errWithDetails.Retry)
	actualDetails := &api.Error{}
	require.NoError(t, anypb.UnmarshalTo(errWithDetails.Details, actualDetails, proto.UnmarshalOptions{}))
	require.True(t, proto.Equal(details, actualDetails), "unexpected details created by WithDetails")
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
