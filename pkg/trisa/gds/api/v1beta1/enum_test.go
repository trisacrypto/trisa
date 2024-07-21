package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
)

func TestEnumParsing(t *testing.T) {
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
