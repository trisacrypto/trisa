package traddr_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp/traddr"
)

func TestEncoding(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		err      error
		msg      string
	}{
		{
			"beneficiary.com/x/12345?t=i",
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			nil,
			"happy case",
		},
	}

	for i, tc := range testCases {
		addr, err := traddr.Encode(tc.input)
		if tc.err != nil {
			require.Error(t, err, "expected error on test case %d: %s", i, tc.msg)
			require.ErrorIs(t, err, tc.err, "error did not match in test case %d: %s", i, tc.msg)
		} else {
			require.NoError(t, err, "expected no error on test case %d: %s", i, tc.msg)
			require.Equal(t, tc.expected, addr, "equality mismatch on test case %d: %s", i, tc.msg)
		}
	}
}

func TestDecoding(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		err      error
		msg      string
	}{
		{
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			"beneficiary.com/x/12345?t=i",
			nil,
			"happy case",
		},
		{
			"2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv", "",
			traddr.ErrUnhandledScheme,
			"error due to missing ta prefix",
		},
	}

	for i, tc := range testCases {
		addr, err := traddr.Decode(tc.input)
		if tc.err != nil {
			require.Error(t, err, "expected error on test case %d: %s", i, tc.msg)
			require.ErrorIs(t, err, tc.err, "error did not match in test case %d: %s", i, tc.msg)
		} else {
			require.NoError(t, err, "expected no error on test case %d: %s", i, tc.msg)
			require.Equal(t, tc.expected, addr, "equality mismatch on test case %d: %s", i, tc.msg)
		}
	}
}
