package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
)

func TestEnumParsing(t *testing.T) {
	t.Run("TransferState", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected api.TransferState
		}{
			{"UNSPECIFIED", api.TransferStateUnspecified},
			{"STARTED", api.TransferStarted},
			{"PENDING", api.TransferPending},
			{"REPAIR", api.TransferRepair},
			{"REVIEW", api.TransferReview},
			{"ACCEPTED", api.TransferAccepted},
			{"COMPLETED", api.TransferCompleted},
			{"REJECTED", api.TransferRejected},

			{"unspecified", api.TransferStateUnspecified},
			{"Started", api.TransferStarted},
			{"PenDinG", api.TransferPending},
			{"RePAIR", api.TransferRepair},
			{"review", api.TransferReview},
			{"aCCepted", api.TransferAccepted},
			{"complETED", api.TransferCompleted},
			{"Rejected", api.TransferRejected},

			{int32(0), api.TransferStateUnspecified},
			{int32(1), api.TransferStarted},
			{int32(2), api.TransferPending},
			{int32(4), api.TransferRepair},
			{int32(3), api.TransferReview},
			{int32(5), api.TransferAccepted},
			{int32(6), api.TransferCompleted},
			{int32(7), api.TransferRejected},
		}

		for i, tc := range validTestCases {
			actual, err := api.ParseTransferState(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"foo", api.ErrUnknownTransferState},
			{2, api.ErrUnknownTransferState},
			{"", api.ErrUnknownTransferState},
			{"TransferState_FOO", api.ErrUnknownTransferState},
			{nil, api.ErrUnknownTransferState},
			{int32(28), api.ErrUnknownTransferState},
		}

		for i, tc := range invalidTestCases {
			actual, err := api.ParseTransferState(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("ConfirmationType", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected api.ConfirmationType
		}{
			{"UNKNOWN", api.ConfirmationTypeUnknown},
			{"SIMPLE", api.ConfirmationTypeSimple},
			{"KEYTOKEN", api.ConfirmationTypeKeyToken},
			{"ONCHAIN", api.ConfirmationTypeOnChain},
			{"unknown", api.ConfirmationTypeUnknown},
			{"Simple", api.ConfirmationTypeSimple},
			{"KeyToken", api.ConfirmationTypeKeyToken},
			{"onCHAIn", api.ConfirmationTypeOnChain},
			{int32(0), api.ConfirmationTypeUnknown},
			{int32(1), api.ConfirmationTypeSimple},
			{int32(2), api.ConfirmationTypeKeyToken},
			{int32(3), api.ConfirmationTypeOnChain},
		}

		for i, tc := range validTestCases {
			actual, err := api.ParseConfirmationType(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"foo", api.ErrUnknownConfirmationType},
			{2, api.ErrUnknownConfirmationType},
			{"", api.ErrUnknownConfirmationType},
			{"ConfirmationType_FOO", api.ErrUnknownConfirmationType},
			{nil, api.ErrUnknownConfirmationType},
			{int32(28), api.ErrUnknownConfirmationType},
		}

		for i, tc := range invalidTestCases {
			actual, err := api.ParseConfirmationType(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("ServiceState", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected api.ServiceState_Status
		}{
			{"UNKNOWN", api.ServiceStatusUnknown},
			{"HEALTHY", api.ServiceStatusHealthy},
			{"UNHEALTHY", api.ServiceStatusUnhealthy},
			{"DANGER", api.ServiceStatusDanger},
			{"OFFLINE", api.ServiceStatusOffline},
			{"MAINTENANCE", api.ServiceStatusMaintenance},
			{"unknown", api.ServiceStatusUnknown},
			{"Healthy", api.ServiceStatusHealthy},
			{"UNhealthy", api.ServiceStatusUnhealthy},
			{"DangeR", api.ServiceStatusDanger},
			{"OFFline", api.ServiceStatusOffline},
			{"maintenance", api.ServiceStatusMaintenance},
			{int32(0), api.ServiceStatusUnknown},
			{int32(1), api.ServiceStatusHealthy},
			{int32(2), api.ServiceStatusUnhealthy},
			{int32(3), api.ServiceStatusDanger},
			{int32(4), api.ServiceStatusOffline},
			{int32(5), api.ServiceStatusMaintenance},
		}

		for i, tc := range validTestCases {
			actual, err := api.ParseServiceState(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"foo", api.ErrUnknownServiceState},
			{2, api.ErrUnknownServiceState},
			{"", api.ErrUnknownServiceState},
			{"ServiceState_FOO", api.ErrUnknownServiceState},
			{nil, api.ErrUnknownServiceState},
			{int32(28), api.ErrUnknownServiceState},
		}

		for i, tc := range invalidTestCases {
			actual, err := api.ParseServiceState(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

}
