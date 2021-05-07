package api_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	protocol "github.com/trisacrypto/trisa/pkg/trisa/api/v1beta1"
)

func TestErrors(t *testing.T) {
	err := protocol.Errorf(protocol.UnknownIdentity, "could not parse %q", "foo")
	require.Error(t, err)
	require.Equal(t, err.Error(), `trisa error [UNKOWN_IDENTITY]: could not parse "foo"`)

	oerr, ok := protocol.Errorp(err)
	require.True(t, ok)
	require.Equal(t, err, oerr)

	oerr, ok = protocol.Errorp(errors.New("unhandled error"))
	require.False(t, ok)
	require.Equal(t, oerr.Error(), "trisa error [UNHANDLED]: unhandled error")

	// TODO: This doesn't work for some reason - the tests just hang?
	// sterr := err.Err()
	// require.Equal(t, sterr.Error(), `rpc error: code = Aborted desc = [UNKOWN_IDENTITY] could not parse "foo"`)

	// oerr, ok = protocol.Errorp(sterr)
	// require.True(t, ok)
	// require.Equal(t, err, oerr)
}
