package trp_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

func TestConfirmationValidate(t *testing.T) {
	testCases := []struct {
		confirm *Confirmation
		err     error
	}{
		{&Confirmation{}, ErrEmptyConfirmation},
		{&Confirmation{Info: &Info{APIVersion: "3.1.0"}}, ErrEmptyConfirmation},
		{&Confirmation{TXID: "foo", Canceled: "bar"}, ErrAmbiguousConfirmation},
		{&Confirmation{TXID: "foo"}, nil},
		{&Confirmation{Canceled: "bar"}, nil},
	}

	for i, tc := range testCases {
		err := tc.confirm.Validate()
		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "test case %d failed with mismatched error", i)
		} else {
			require.NoError(t, err, "test case %d failed: expected valid confirmation", i)
		}
	}
}
