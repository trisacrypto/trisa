package envelope_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/trisa/envelope"
)

func TestStateType(t *testing.T) {
	require.Equal(t, "clear", envelope.Clear.String())
	require.Equal(t, "unsealed", envelope.Unsealed.String())
	require.Equal(t, "sealed", envelope.Sealed.String())
	require.Equal(t, "error", envelope.Error.String())
	require.Equal(t, "clear-error", envelope.ClearError.String())
	require.Equal(t, "unsealed-error", envelope.UnsealedError.String())
	require.Equal(t, "sealed-error", envelope.SealedError.String())
	require.Equal(t, "corrupted", envelope.Corrupted.String())

	require.Equal(t, "unknown", envelope.State(0).String())
	require.Equal(t, "unknown", envelope.State(42).String())

	// These only work because of the current positions of the states, e.g. Error==4
	// and Clear-Sealed and ClearError-UnsealedError are 1-3 and 5-7 respectively.
	require.Equal(t, envelope.ClearError, envelope.Clear|envelope.Error)
	require.Equal(t, envelope.UnsealedError, envelope.Unsealed|envelope.Error)
	require.Equal(t, envelope.SealedError, envelope.Sealed|envelope.Error)
}

func TestEnvelopeState(t *testing.T) {
	// No payload, no envelope should return an invalid state
	env := &envelope.Envelope{}
	require.Equal(t, envelope.Unknown, env.State())
}
