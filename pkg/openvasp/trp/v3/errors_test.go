package trp_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/trisa/pkg/openvasp/trp/v3"
)

func TestStatusError(t *testing.T) {
	testCases := []struct {
		err      error
		expected string
	}{
		{&StatusError{http.StatusBadRequest, "something bad happened"}, "something bad happened"},
		{&StatusError{http.StatusBadRequest, ""}, "Bad Request"},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, tc.err.Error(), "test case %d failed", i)
	}
}
