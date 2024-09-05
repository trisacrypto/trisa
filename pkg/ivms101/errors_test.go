package ivms101_test

import (
	"encoding/json"
	"errors"
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

func TestValidationErrorHelpers(t *testing.T) {
	tests := []struct {
		err      error
		expected string
	}{
		{
			ivms101.MissingField("name"),
			"ivms101: missing name: this field is required",
		},
		{
			ivms101.IncorrectField("country", "must be a valid iso 3166-1 alpha-2 code"),
			"ivms101: invalid field country: must be a valid iso 3166-1 alpha-2 code",
		},
		{
			ivms101.MaxNText("customerIdentification", 50, 63),
			"ivms101: customerIdentification exceeded max length 50 chars: 63 characters is too long",
		},
		{
			ivms101.OneOfMissing("naturalPerson", "legalPerson"),
			"ivms101: missing one of naturalPerson or legalPerson: at most one of these fields is required",
		},
		{
			ivms101.OneOfMissing("naturalPerson"),
			"ivms101: missing naturalPerson: this field is required",
		},
		{
			ivms101.OneOfMissing("naturalPerson", "legalPerson", "animalPerson", "spiritPerson"),
			"ivms101: missing one of naturalPerson, legalPerson, animalPerson, or spiritPerson: at most one of these fields is required",
		},
		{
			ivms101.OneOfTooMany("naturalPerson", "legalPerson"),
			"ivms101: specify only one of naturalPerson or legalPerson: at most one of these fields may be specified",
		},
		{
			ivms101.OneOfTooMany("naturalPerson", "legalPerson", "animalPerson", "spiritPerson"),
			"ivms101: specify only one of naturalPerson, legalPerson, animalPerson, or spiritPerson: at most one of these fields may be specified",
		},
	}

	for i, tc := range tests {
		require.EqualError(t, tc.err, tc.expected, "test %d failed", i)
	}
}

func TestValidationErrorPanic(t *testing.T) {
	require.Panics(t, func() {
		ivms101.OneOfMissing()
	})

	require.Panics(t, func() {
		ivms101.OneOfTooMany()
	})

	require.Panics(t, func() {
		ivms101.OneOfTooMany("foo")
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		tests := []error{
			ivms101.ValidationError("", nil),
			ivms101.ValidationError("", nil, nil, nil, nil, nil),
			ivms101.ValidationError("bang", nil, nil, nil, nil, nil),
		}

		for i, err := range tests {
			require.NoError(t, err, "test case %d failed", i)
		}
	})

	t.Run("Single", func(t *testing.T) {
		tests := []struct {
			err      error
			expected string
		}{
			{
				ivms101.ValidationError("", nil, ivms101.MissingField("foo")),
				"ivms101: missing foo: this field is required",
			},
			{
				ivms101.ValidationError("", ivms101.ValidationError("", ivms101.MissingField("foo"), nil), nil),
				"ivms101: missing foo: this field is required",
			},
			{
				ivms101.ValidationError("", nil, errors.New("something bad happened")),
				"ivms101: invalid input: something bad happened",
			},
			{
				ivms101.ValidationError("", errors.New("something bad happened"), nil),
				"ivms101: invalid input: something bad happened",
			},
			{
				ivms101.ValidationError("", nil, ivms101.ValidationError("", ivms101.MissingField("foo"), nil)),
				"ivms101: missing foo: this field is required",
			},
		}

		for i, tc := range tests {
			require.EqualError(t, tc.err, tc.expected, "test case %d failed", i)
		}
	})

	t.Run("Multi", func(t *testing.T) {
		tests := []struct {
			err      error
			expected string
		}{
			{
				ivms101.ValidationError("", nil, ivms101.MissingField("foo"), ivms101.IncorrectField("bar", "bad bar")),
				"2 validation errors occurred:\n  ivms101: missing foo: this field is required\n  ivms101: invalid field bar: bad bar",
			},
			{
				ivms101.ValidationError("", ivms101.MissingField("name"), ivms101.MissingField("foo"), ivms101.IncorrectField("bar", "bad bar")),
				"3 validation errors occurred:\n  ivms101: missing name: this field is required\n  ivms101: missing foo: this field is required\n  ivms101: invalid field bar: bad bar",
			},
			{
				ivms101.ValidationError("", ivms101.MissingField("name"), ivms101.ValidationError("parent", nil, ivms101.ValidationError("sub", nil, ivms101.MissingField("foo")), ivms101.IncorrectField("bar", "bad bar"))),
				"3 validation errors occurred:\n  ivms101: missing name: this field is required\n  ivms101: missing parent.sub.foo: this field is required\n  ivms101: invalid field parent.bar: bad bar",
			},
			{
				ivms101.ValidationError("super", ivms101.MissingField("name"), ivms101.ValidationError("parent", ivms101.ValidationError("sub", nil, ivms101.MissingField("foo")), ivms101.IncorrectField("bar", "bad bar"))),
				"3 validation errors occurred:\n  ivms101: missing name: this field is required\n  ivms101: missing super.sub.foo: this field is required\n  ivms101: invalid field super.parent.bar: bad bar",
			},
		}

		for i, tc := range tests {
			require.EqualError(t, tc.err, tc.expected, "test case %d failed", i)
		}
	})
}
