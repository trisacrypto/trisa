package ivms101_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func TestWrap(t *testing.T) {
	cause := &json.UnsupportedValueError{Str: "this value is unsupported"}
	err := ivms101.Wrap(ivms101.ErrInvalidNaturalPersonNameTypeCode, cause)

	require.ErrorIs(t, err, ivms101.ErrInvalidNaturalPersonNameTypeCode)
	require.ErrorIs(t, err, cause)
	require.EqualError(t, err, ivms101.ErrInvalidNaturalPersonNameTypeCode.Error())
}
