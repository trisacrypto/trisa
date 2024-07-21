package models_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	models "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
)

func TestEnumParsing(t *testing.T) {
	t.Run("BusinessCategory", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected models.BusinessCategory
		}{
			// Original Tests
			{"unknown entity", models.BusinessCategoryUnknown},
			{"PRIVATE_ORGANIZATION", models.BusinessCategoryPrivate},
			{"Government Entity", models.BusinessCategoryGovernment},
			{"Business_entity", models.BusinessCategoryBusiness},
			{"non commercial entity", models.BusinessCategoryNonCommercial},

			{int32(0), models.BusinessCategoryUnknown},
			{int32(1), models.BusinessCategoryPrivate},
			{int32(2), models.BusinessCategoryGovernment},
			{int32(3), models.BusinessCategoryBusiness},
			{int32(4), models.BusinessCategoryNonCommercial},
		}

		for i, tc := range validTestCases {
			actual, err := models.ParseBusinessCategory(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"foo bar", models.ErrUnknownBusinessCategory},
			{2, models.ErrUnknownBusinessCategory},
			{"", models.ErrUnknownBusinessCategory},
			{"BusinessCategory_FOO", models.ErrUnknownBusinessCategory},
			{nil, models.ErrUnknownBusinessCategory},
			{int32(28), models.ErrUnknownBusinessCategory},
		}

		for i, tc := range invalidTestCases {
			actual, err := models.ParseBusinessCategory(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("VASPCategory", func(t *testing.T) {
		validTestCases := []struct {
			input    string
			expected string
		}{
			{"Unknown", models.VASPCategoryUnknown},
			{"Exchange", models.VASPCategoryExchange},
			{"DEX", models.VASPCategoryDEX},
			{"P2P", models.VASPCategoryP2P},
			{"Kiosk", models.VASPCategoryKiosk},
			{"Custodian", models.VASPCategoryCustodian},
			{"OTC", models.VASPCategoryOTC},
			{"Fund", models.VASPCategoryFund},
			{"Project", models.VASPCategoryProject},
			{"Gambling", models.VASPCategoryGambling},
			{"Miner", models.VASPCategoryMiner},
			{"Mixer", models.VASPCategoryMixer},
			{"Individual", models.VASPCategoryIndividual},
			{"Other", models.VASPCategoryOther},
			{"unknown", models.VASPCategoryUnknown},
			{"EXCHANGE", models.VASPCategoryExchange},
			{"dex", models.VASPCategoryDEX},
			{"P2p", models.VASPCategoryP2P},
			{"KiosK", models.VASPCategoryKiosk},
			{"CustODian", models.VASPCategoryCustodian},
			{"otc", models.VASPCategoryOTC},
			{"FuNd", models.VASPCategoryFund},
			{"ProJeCt", models.VASPCategoryProject},
			{"GamblinG", models.VASPCategoryGambling},
			{"MineR", models.VASPCategoryMiner},
			{"MiXer", models.VASPCategoryMixer},
			{"individual", models.VASPCategoryIndividual},
			{"other", models.VASPCategoryOther},
		}

		for i, tc := range validTestCases {
			actual, err := models.ValidVASPCategory(tc.input)
			require.NoError(t, err, "unexpected error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "unexpected value returned for valid test case %d", i)
		}

		invalidTestCases := []struct {
			input string
			err   error
		}{
			{"p3p", models.ErrUnknownVASPCategory},
			{"", models.ErrUnknownVASPCategory},
			{"   Project   ", models.ErrUnknownVASPCategory},
		}

		for i, tc := range invalidTestCases {
			actual, err := models.ValidVASPCategory(tc.input)
			require.ErrorIs(t, err, tc.err, "expected an error on invalid test case %d", i)
			require.Empty(t, actual, "expected empty string returned on invalid test case %d", i)
		}
	})

	t.Run("VerificationState", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected models.VerificationState
		}{
			{"NO_VERIFICATION", models.VerificationNone},
			{"SUBMITTED", models.VerificationSubmitted},
			{"EMAIL_VERIFIED", models.VerificationEmailVerified},
			{"PENDING_REVIEW", models.VerificationPending},
			{"REVIEWED", models.VerificationReviewed},
			{"ISSUING_CERTIFICATE", models.VerificationIssuing},
			{"VERIFIED", models.VerificationVerified},
			{"REJECTED", models.VerificationRejected},
			{"APPEALED", models.VerificationAppealed},
			{"ERRORED", models.VerificationErrored},

			{"No Verification", models.VerificationNone},
			{"submitted", models.VerificationSubmitted},
			{"email_VERIFIED", models.VerificationEmailVerified},
			{"pending review", models.VerificationPending},
			{"reviewed", models.VerificationReviewed},
			{"Issuing CERTIFICATE", models.VerificationIssuing},
			{"Verified", models.VerificationVerified},
			{"RejecteD", models.VerificationRejected},
			{"APPealed", models.VerificationAppealed},
			{"errored", models.VerificationErrored},

			{int32(0), models.VerificationNone},
			{int32(1), models.VerificationSubmitted},
			{int32(2), models.VerificationEmailVerified},
			{int32(3), models.VerificationPending},
			{int32(4), models.VerificationReviewed},
			{int32(5), models.VerificationIssuing},
			{int32(6), models.VerificationVerified},
			{int32(7), models.VerificationRejected},
			{int32(8), models.VerificationAppealed},
			{int32(9), models.VerificationErrored},
		}

		for i, tc := range validTestCases {
			actual, err := models.ParseVerificationState(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"foo", models.ErrUnknownVerificationState},
			{2, models.ErrUnknownVerificationState},
			{"", models.ErrUnknownVerificationState},
			{"ServiceState_FOO", models.ErrUnknownVerificationState},
			{nil, models.ErrUnknownVerificationState},
			{int32(28), models.ErrUnknownVerificationState},
		}

		for i, tc := range invalidTestCases {
			actual, err := models.ParseVerificationState(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("ServiceState", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected models.ServiceState
		}{
			{"UNKNOWN", models.ServiceStatusUnknown},
			{"HEALTHY", models.ServiceStatusHealthy},
			{"UNHEALTHY", models.ServiceStatusUnhealthy},
			{"DANGER", models.ServiceStatusDanger},
			{"OFFLINE", models.ServiceStatusOffline},
			{"MAINTENANCE", models.ServiceStatusMaintenance},
			{"unknown", models.ServiceStatusUnknown},
			{"Healthy", models.ServiceStatusHealthy},
			{"UNhealthy", models.ServiceStatusUnhealthy},
			{"DangeR", models.ServiceStatusDanger},
			{"OFFline", models.ServiceStatusOffline},
			{"maintenance", models.ServiceStatusMaintenance},
			{int32(0), models.ServiceStatusUnknown},
			{int32(1), models.ServiceStatusHealthy},
			{int32(2), models.ServiceStatusUnhealthy},
			{int32(3), models.ServiceStatusDanger},
			{int32(4), models.ServiceStatusOffline},
			{int32(5), models.ServiceStatusMaintenance},
		}

		for i, tc := range validTestCases {
			actual, err := models.ParseServiceState(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"foo", models.ErrUnknownServiceState},
			{2, models.ErrUnknownServiceState},
			{"", models.ErrUnknownServiceState},
			{"ServiceState_FOO", models.ErrUnknownServiceState},
			{nil, models.ErrUnknownServiceState},
			{int32(28), models.ErrUnknownServiceState},
		}

		for i, tc := range invalidTestCases {
			actual, err := models.ParseServiceState(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})
}
