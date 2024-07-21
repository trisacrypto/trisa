package ivms101_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func TestEnumParse(t *testing.T) {
	t.Run("NaturalPersonNameTypeCode", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected ivms101.NaturalPersonNameTypeCode
		}{
			{"ALIA", ivms101.NaturalPersonAlias},
			{"BIRT", ivms101.NaturalPersonBirth},
			{"MAID", ivms101.NaturalPersonMaiden},
			{"LEGL", ivms101.NaturalPersonLegal},
			{"MISC", ivms101.NaturalPersonMisc},
			{"NATURAL_PERSON_NAME_TYPE_CODE_ALIA", ivms101.NaturalPersonAlias},
			{"NATURAL_PERSON_NAME_TYPE_CODE_BIRT", ivms101.NaturalPersonBirth},
			{"NATURAL_PERSON_NAME_TYPE_CODE_MAID", ivms101.NaturalPersonMaiden},
			{"NATURAL_PERSON_NAME_TYPE_CODE_LEGL", ivms101.NaturalPersonLegal},
			{"NATURAL_PERSON_NAME_TYPE_CODE_MISC", ivms101.NaturalPersonMisc},
			{"alia", ivms101.NaturalPersonAlias},
			{"Birt", ivms101.NaturalPersonBirth},
			{"mAiD", ivms101.NaturalPersonMaiden},
			{"LegL", ivms101.NaturalPersonLegal},
			{"misc", ivms101.NaturalPersonMisc},
			{int32(1), ivms101.NaturalPersonAlias},
			{int32(2), ivms101.NaturalPersonBirth},
			{int32(3), ivms101.NaturalPersonMaiden},
			{int32(4), ivms101.NaturalPersonLegal},
			{int32(0), ivms101.NaturalPersonMisc},
		}

		for i, tc := range validTestCases {
			actual, err := ivms101.ParseNaturalPersonNameTypeCode(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"alias", ivms101.ErrCouldNotParseEnum},
			{2, ivms101.ErrCouldNotParseEnum},
			{"", ivms101.ErrCouldNotParseEnum},
			{"NATURAL_PERSON_NAME_TYPE_CODE_FOO", ivms101.ErrCouldNotParseEnum},
			{nil, ivms101.ErrCouldNotParseEnum},
			{int32(28), ivms101.ErrCouldNotParseEnum},
		}

		for i, tc := range invalidTestCases {
			actual, err := ivms101.ParseNaturalPersonNameTypeCode(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("LegalPersonNameTypeCode", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected ivms101.LegalPersonNameTypeCode
		}{
			{"MISC", ivms101.LegalPersonMisc},
			{"LEGL", ivms101.LegalPersonLegal},
			{"SHRT", ivms101.LegalPersonShort},
			{"TRAD", ivms101.LegalPersonTrading},
			{"LEGAL_PERSON_NAME_TYPE_CODE_MISC", ivms101.LegalPersonMisc},
			{"LEGAL_PERSON_NAME_TYPE_CODE_LEGL", ivms101.LegalPersonLegal},
			{"LEGAL_PERSON_NAME_TYPE_CODE_SHRT", ivms101.LegalPersonShort},
			{"LEGAL_PERSON_NAME_TYPE_CODE_TRAD", ivms101.LegalPersonTrading},
			{"legl", ivms101.LegalPersonLegal},
			{"Shrt", ivms101.LegalPersonShort},
			{"TraD", ivms101.LegalPersonTrading},
			{int32(0), ivms101.LegalPersonMisc},
			{int32(1), ivms101.LegalPersonLegal},
			{int32(2), ivms101.LegalPersonShort},
			{int32(3), ivms101.LegalPersonTrading},
		}

		for i, tc := range validTestCases {
			actual, err := ivms101.ParseLegalPersonNameTypeCode(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"legal", ivms101.ErrCouldNotParseEnum},
			{2, ivms101.ErrCouldNotParseEnum},
			{"", ivms101.ErrCouldNotParseEnum},
			{"LEGAL_PERSON_NAME_TYPE_CODE_FOO", ivms101.ErrCouldNotParseEnum},
			{nil, ivms101.ErrCouldNotParseEnum},
			{int32(28), ivms101.ErrCouldNotParseEnum},
		}

		for i, tc := range invalidTestCases {
			actual, err := ivms101.ParseLegalPersonNameTypeCode(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("AddressTypeCode", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected ivms101.AddressTypeCode
		}{
			{"MISC", ivms101.AddressTypeMisc},
			{"HOME", ivms101.AddressTypeHome},
			{"BIZZ", ivms101.AddressTypeBusiness},
			{"GEOG", ivms101.AddressTypeGeographic},
			{"ADDRESS_TYPE_CODE_MISC", ivms101.AddressTypeMisc},
			{"ADDRESS_TYPE_CODE_HOME", ivms101.AddressTypeHome},
			{"ADDRESS_TYPE_CODE_BIZZ", ivms101.AddressTypeBusiness},
			{"ADDRESS_TYPE_CODE_GEOG", ivms101.AddressTypeGeographic},
			{"misc", ivms101.AddressTypeMisc},
			{"Home", ivms101.AddressTypeHome},
			{"bIzZ", ivms101.AddressTypeBusiness},
			{"GeoG", ivms101.AddressTypeGeographic},
			{int32(0), ivms101.AddressTypeMisc},
			{int32(1), ivms101.AddressTypeHome},
			{int32(2), ivms101.AddressTypeBusiness},
			{int32(3), ivms101.AddressTypeGeographic},
		}

		for i, tc := range validTestCases {
			actual, err := ivms101.ParseAddressTypeCode(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"business", ivms101.ErrCouldNotParseEnum},
			{2, ivms101.ErrCouldNotParseEnum},
			{"", ivms101.ErrCouldNotParseEnum},
			{"ADDRESS_TYPE_CODE_FOO", ivms101.ErrCouldNotParseEnum},
			{nil, ivms101.ErrCouldNotParseEnum},
			{int32(28), ivms101.ErrCouldNotParseEnum},
		}

		for i, tc := range invalidTestCases {
			actual, err := ivms101.ParseAddressTypeCode(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("NationalIdentifierTypeCode", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected ivms101.NationalIdentifierTypeCode
		}{
			{"ARNU", ivms101.NationalIdentifierARNU},
			{"CCPT", ivms101.NationalIdentifierCCPT},
			{"RAID", ivms101.NationalIdentifierRAID},
			{"DRLC", ivms101.NationalIdentifierDRLC},
			{"FIIN", ivms101.NationalIdentifierFIIN},
			{"TXID", ivms101.NationalIdentifierTXID},
			{"SOCS", ivms101.NationalIdentifierSOCS},
			{"IDCD", ivms101.NationalIdentifierIDCD},
			{"LEIX", ivms101.NationalIdentifierLEIX},
			{"MISC", ivms101.NationalIdentifierMISC},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_ARNU", ivms101.NationalIdentifierARNU},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_CCPT", ivms101.NationalIdentifierCCPT},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_RAID", ivms101.NationalIdentifierRAID},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_DRLC", ivms101.NationalIdentifierDRLC},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_FIIN", ivms101.NationalIdentifierFIIN},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_TXID", ivms101.NationalIdentifierTXID},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_SOCS", ivms101.NationalIdentifierSOCS},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_IDCD", ivms101.NationalIdentifierIDCD},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_LEIX", ivms101.NationalIdentifierLEIX},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_MISC", ivms101.NationalIdentifierMISC},
			{"arnu", ivms101.NationalIdentifierARNU},
			{"CcPt", ivms101.NationalIdentifierCCPT},
			{"Raid", ivms101.NationalIdentifierRAID},
			{"drlc", ivms101.NationalIdentifierDRLC},
			{"fIIn", ivms101.NationalIdentifierFIIN},
			{"TxiD", ivms101.NationalIdentifierTXID},
			{"sOCs", ivms101.NationalIdentifierSOCS},
			{"Idcd", ivms101.NationalIdentifierIDCD},
			{"leIX", ivms101.NationalIdentifierLEIX},
			{"Misc", ivms101.NationalIdentifierMISC},
			{int32(1), ivms101.NationalIdentifierARNU},
			{int32(2), ivms101.NationalIdentifierCCPT},
			{int32(3), ivms101.NationalIdentifierRAID},
			{int32(4), ivms101.NationalIdentifierDRLC},
			{int32(5), ivms101.NationalIdentifierFIIN},
			{int32(6), ivms101.NationalIdentifierTXID},
			{int32(7), ivms101.NationalIdentifierSOCS},
			{int32(8), ivms101.NationalIdentifierIDCD},
			{int32(9), ivms101.NationalIdentifierLEIX},
			{int32(0), ivms101.NationalIdentifierMISC},
		}

		for i, tc := range validTestCases {
			actual, err := ivms101.ParseNationalIdentifierTypeCode(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"passport", ivms101.ErrCouldNotParseEnum},
			{2, ivms101.ErrCouldNotParseEnum},
			{"", ivms101.ErrCouldNotParseEnum},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_FOO", ivms101.ErrCouldNotParseEnum},
			{nil, ivms101.ErrCouldNotParseEnum},
			{int32(28), ivms101.ErrCouldNotParseEnum},
		}

		for i, tc := range invalidTestCases {
			actual, err := ivms101.ParseNationalIdentifierTypeCode(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})

	t.Run("TransliterationMethodCode", func(t *testing.T) {
		validTestCases := []struct {
			input    any
			expected ivms101.TransliterationMethodCode
		}{
			{"OTHR", ivms101.TransliterationMethodOTHR},
			{"ARAB", ivms101.TransliterationMethodARAB},
			{"ARAN", ivms101.TransliterationMethodARAN},
			{"ARMN", ivms101.TransliterationMethodARMN},
			{"CYRL", ivms101.TransliterationMethodCYRL},
			{"DEVA", ivms101.TransliterationMethodDEVA},
			{"GEOR", ivms101.TransliterationMethodGEOR},
			{"GREK", ivms101.TransliterationMethodGREK},
			{"HANI", ivms101.TransliterationMethodHANI},
			{"HEBR", ivms101.TransliterationMethodHEBR},
			{"KANA", ivms101.TransliterationMethodKANA},
			{"KORE", ivms101.TransliterationMethodKORE},
			{"THAI", ivms101.TransliterationMethodTHAI},
			{"TRANSLITERATION_METHOD_CODE_OTHR", ivms101.TransliterationMethodOTHR},
			{"TRANSLITERATION_METHOD_CODE_ARAB", ivms101.TransliterationMethodARAB},
			{"TRANSLITERATION_METHOD_CODE_ARAN", ivms101.TransliterationMethodARAN},
			{"TRANSLITERATION_METHOD_CODE_ARMN", ivms101.TransliterationMethodARMN},
			{"TRANSLITERATION_METHOD_CODE_CYRL", ivms101.TransliterationMethodCYRL},
			{"TRANSLITERATION_METHOD_CODE_DEVA", ivms101.TransliterationMethodDEVA},
			{"TRANSLITERATION_METHOD_CODE_GEOR", ivms101.TransliterationMethodGEOR},
			{"TRANSLITERATION_METHOD_CODE_GREK", ivms101.TransliterationMethodGREK},
			{"TRANSLITERATION_METHOD_CODE_HANI", ivms101.TransliterationMethodHANI},
			{"TRANSLITERATION_METHOD_CODE_HEBR", ivms101.TransliterationMethodHEBR},
			{"TRANSLITERATION_METHOD_CODE_KANA", ivms101.TransliterationMethodKANA},
			{"TRANSLITERATION_METHOD_CODE_KORE", ivms101.TransliterationMethodKORE},
			{"TRANSLITERATION_METHOD_CODE_THAI", ivms101.TransliterationMethodTHAI},
			{"othr", ivms101.TransliterationMethodOTHR},
			{"Arab", ivms101.TransliterationMethodARAB},
			{"ARan", ivms101.TransliterationMethodARAN},
			{"ArmN", ivms101.TransliterationMethodARMN},
			{"cYRl", ivms101.TransliterationMethodCYRL},
			{"devA", ivms101.TransliterationMethodDEVA},
			{"GeOR", ivms101.TransliterationMethodGEOR},
			{"GreK", ivms101.TransliterationMethodGREK},
			{"Hani", ivms101.TransliterationMethodHANI},
			{"hebr", ivms101.TransliterationMethodHEBR},
			{"kAnA", ivms101.TransliterationMethodKANA},
			{"Kore", ivms101.TransliterationMethodKORE},
			{"Thai", ivms101.TransliterationMethodTHAI},
			{int32(0), ivms101.TransliterationMethodOTHR},
			{int32(1), ivms101.TransliterationMethodARAB},
			{int32(2), ivms101.TransliterationMethodARAN},
			{int32(3), ivms101.TransliterationMethodARMN},
			{int32(4), ivms101.TransliterationMethodCYRL},
			{int32(5), ivms101.TransliterationMethodDEVA},
			{int32(6), ivms101.TransliterationMethodGEOR},
			{int32(7), ivms101.TransliterationMethodGREK},
			{int32(8), ivms101.TransliterationMethodHANI},
			{int32(9), ivms101.TransliterationMethodHEBR},
			{int32(10), ivms101.TransliterationMethodKANA},
			{int32(11), ivms101.TransliterationMethodKORE},
			{int32(12), ivms101.TransliterationMethodTHAI},
		}

		for i, tc := range validTestCases {
			actual, err := ivms101.ParseTransliterationMethodCode(tc.input)
			require.NoError(t, err, "did not expect error on valid test case %d", i)
			require.Equal(t, tc.expected, actual, "mismatched expectation on valid test case %d", i)
		}

		invalidTestCases := []struct {
			input any
			err   error
		}{
			{"arabic", ivms101.ErrCouldNotParseEnum},
			{2, ivms101.ErrCouldNotParseEnum},
			{"", ivms101.ErrCouldNotParseEnum},
			{"TRANSLITERATION_METHOD_CODE_FOO", ivms101.ErrCouldNotParseEnum},
			{nil, ivms101.ErrCouldNotParseEnum},
			{int32(28), ivms101.ErrCouldNotParseEnum},
		}

		for i, tc := range invalidTestCases {
			actual, err := ivms101.ParseTransliterationMethodCode(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})
}

func TestMarshalNatNameCode(t *testing.T) {
	// Test an array of expected ivms101 protocol buffer enums
	// to inspect whether we correctly marshal them to the correct
	// compatible ivms101 format
	expected := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"ALIA"`),
		2: []byte(`"BIRT"`),
		3: []byte(`"MAID"`),
		4: []byte(`"LEGL"`),
	}
	for c := 0; c < 5; c++ {
		code := ivms101.NaturalPersonNameTypeCode(c)
		data, err := code.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, expected[c], data)
	}
}

func TestUnmarshalNatNameCode(t *testing.T) {
	// Test an array of expected compatible ivms101 codes
	// to inspect whether we correctly unmarshal them to the correct
	// protocol buffer enum
	data := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"ALIA"`),
		2: []byte(`"BIRT"`),
		3: []byte(`"MAID"`),
		4: []byte(`"LEGL"`),
	}
	for c := 0; c < 5; c++ {
		var code ivms101.NaturalPersonNameTypeCode
		err := code.UnmarshalJSON(data[c])
		require.Nil(t, err)
		require.Equal(t, ivms101.NaturalPersonNameTypeCode(c), code)
	}

	// Test than an unknown compatible ivms101 code triggers a helpful error
	var unknown ivms101.NaturalPersonNameTypeCode
	err := unknown.UnmarshalJSON([]byte(`"STAGE"`))
	require.EqualError(t, err, "invalid NaturalPersonNameTypeCode alias")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.NaturalPersonNameTypeCode(0))

	// Test that incorrect json input (no quotes) triggers a helpful error
	var badCode ivms101.NaturalPersonNameTypeCode
	err = badCode.UnmarshalJSON([]byte("ALIA"))
	require.EqualError(t, err, "could not parse NaturalPersonNameTypeCode from value")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.NaturalPersonNameTypeCode(0))
}

func TestMarshalLegNameCode(t *testing.T) {
	// Test an array of expected ivms101 protocol buffer enums
	// to inspect whether we correctly marshal them to the correct
	// compatible ivms101 format
	expected := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"LEGL"`),
		2: []byte(`"SHRT"`),
		3: []byte(`"TRAD"`),
	}
	for c := 0; c < 4; c++ {
		code := ivms101.LegalPersonNameTypeCode(c)
		data, err := code.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, expected[c], data)
	}
}

func TestUnmarshalLegNameCode(t *testing.T) {
	// Test an array of expected compatible ivms101 codes
	// to inspect whether we correctly unmarshal them to the correct
	// protocol buffer enum
	data := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"LEGL"`),
		2: []byte(`"SHRT"`),
		3: []byte(`"TRAD"`),
	}
	for c := 0; c < 4; c++ {
		var code ivms101.LegalPersonNameTypeCode
		err := code.UnmarshalJSON(data[c])
		require.Nil(t, err)
		require.Equal(t, ivms101.LegalPersonNameTypeCode(c), code)
	}

	// Test than an unknown compatible ivms101 code triggers a helpful error
	var unknown ivms101.LegalPersonNameTypeCode
	err := unknown.UnmarshalJSON([]byte(`"SHEL"`))
	require.EqualError(t, err, "invalid LegalPersonNameTypeCode alias")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.LegalPersonNameTypeCode(0))

	// Test that incorrect json input (no quotes) triggers a helpful error
	var badCode ivms101.LegalPersonNameTypeCode
	err = badCode.UnmarshalJSON([]byte("LEGL"))
	require.EqualError(t, err, "could not parse LegalPersonNameTypeCode from value")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.LegalPersonNameTypeCode(0))
}

func TestMarshalAddrCode(t *testing.T) {
	// Test an array of expected ivms101 protocol buffer enums
	// to inspect whether we correctly marshal them to the correct
	// compatible ivms101 format
	expected := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"HOME"`),
		2: []byte(`"BIZZ"`),
		3: []byte(`"GEOG"`),
	}
	for c := 0; c < 4; c++ {
		code := ivms101.AddressTypeCode(c)
		data, err := code.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, expected[c], data)
	}
}

func TestUnmarshalAddrCode(t *testing.T) {
	// Test an array of expected compatible ivms101 codes
	// to inspect whether we correctly unmarshal them to the correct
	// protocol buffer enum
	data := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"HOME"`),
		2: []byte(`"BIZZ"`),
		3: []byte(`"GEOG"`),
	}
	for c := 0; c < 4; c++ {
		var code ivms101.AddressTypeCode
		err := code.UnmarshalJSON(data[c])
		require.Nil(t, err)
		require.Equal(t, ivms101.AddressTypeCode(c), code)
	}

	// Test than an unknown compatible ivms101 code triggers a helpful error
	var unknown ivms101.AddressTypeCode
	err := unknown.UnmarshalJSON([]byte(`"LALA"`))
	require.EqualError(t, err, "invalid AddressTypeCode alias")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.AddressTypeCode(0))

	// Test that incorrect json input (no quotes) triggers a helpful error
	var badCode ivms101.AddressTypeCode
	err = badCode.UnmarshalJSON([]byte("HOME"))
	require.EqualError(t, err, "could not parse AddressTypeCode from value")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.AddressTypeCode(0))
}

func TestMarshalNatIDCode(t *testing.T) {
	// Test an array of expected ivms101 protocol buffer enums
	// to inspect whether we correctly marshal them to the correct
	// compatible ivms101 format
	expected := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"ARNU"`),
		2: []byte(`"CCPT"`),
		3: []byte(`"RAID"`),
		4: []byte(`"DRLC"`),
		5: []byte(`"FIIN"`),
		6: []byte(`"TXID"`),
		7: []byte(`"SOCS"`),
		8: []byte(`"IDCD"`),
		9: []byte(`"LEIX"`),
	}
	for c := 0; c < 10; c++ {
		code := ivms101.NationalIdentifierTypeCode(c)
		data, err := code.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, expected[c], data)
	}
}

func TestUnmarshalNatIDCode(t *testing.T) {
	// Test an array of expected compatible ivms101 codes
	// to inspect whether we correctly unmarshal them to the correct
	// protocol buffer enum
	data := map[int][]byte{
		0: []byte(`"MISC"`),
		1: []byte(`"ARNU"`),
		2: []byte(`"CCPT"`),
		3: []byte(`"RAID"`),
		4: []byte(`"DRLC"`),
		5: []byte(`"FIIN"`),
		6: []byte(`"TXID"`),
		7: []byte(`"SOCS"`),
		8: []byte(`"IDCD"`),
		9: []byte(`"LEIX"`),
	}
	for c := 0; c < 10; c++ {
		var code ivms101.NationalIdentifierTypeCode
		err := code.UnmarshalJSON(data[c])
		require.Nil(t, err)
		require.Equal(t, ivms101.NationalIdentifierTypeCode(c), code)
	}

	// Test than an unknown compatible ivms101 code triggers a helpful error
	var unknown ivms101.NationalIdentifierTypeCode
	err := unknown.UnmarshalJSON([]byte(`"ACME"`))
	require.EqualError(t, err, "invalid NationalIdentifierTypeCode alias")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.NationalIdentifierTypeCode(0))

	// Test that incorrect json input (no quotes) triggers a helpful error
	var badCode ivms101.NationalIdentifierTypeCode
	err = badCode.UnmarshalJSON([]byte("LEIX"))
	require.EqualError(t, err, "could not parse NationalIdentifierTypeCode from value")
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.NationalIdentifierTypeCode(0))
}

func TestMarshalTransliterationCode(t *testing.T) {
	// Test an array of expected ivms101 protocol buffer enums
	// to inspect whether we correctly marshal them to the correct
	// compatible ivms101 format
	expected := map[int][]byte{
		0:  []byte(`"OTHR"`),
		1:  []byte(`"ARAB"`),
		2:  []byte(`"ARAN"`),
		3:  []byte(`"ARMN"`),
		4:  []byte(`"CYRL"`),
		5:  []byte(`"DEVA"`),
		6:  []byte(`"GEOR"`),
		7:  []byte(`"GREK"`),
		8:  []byte(`"HANI"`),
		9:  []byte(`"HEBR"`),
		10: []byte(`"KANA"`),
		11: []byte(`"KORE"`),
		12: []byte(`"THAI"`),
	}
	for c := 0; c < 10; c++ {
		code := ivms101.TransliterationMethodCode(c)
		data, err := code.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, expected[c], data)
	}
}

func TestUnmarshalTransliterationCode(t *testing.T) {
	// Test an array of expected compatible ivms101 codes
	// to inspect whether we correctly unmarshal them to the correct
	// protocol buffer enum
	data := map[int][]byte{
		0:  []byte(`"OTHR"`),
		1:  []byte(`"ARAB"`),
		2:  []byte(`"ARAN"`),
		3:  []byte(`"ARMN"`),
		4:  []byte(`"CYRL"`),
		5:  []byte(`"DEVA"`),
		6:  []byte(`"GEOR"`),
		7:  []byte(`"GREK"`),
		8:  []byte(`"HANI"`),
		9:  []byte(`"HEBR"`),
		10: []byte(`"KANA"`),
		11: []byte(`"KORE"`),
		12: []byte(`"THAI"`),
	}
	for c := 0; c < 10; c++ {
		var code ivms101.TransliterationMethodCode
		err := code.UnmarshalJSON(data[c])
		require.Nil(t, err)
		require.Equal(t, ivms101.TransliterationMethodCode(c), code)
	}

	// Test than an unknown compatible ivms101 code triggers a helpful error
	var unknown ivms101.TransliterationMethodCode
	err := unknown.UnmarshalJSON([]byte(`"KLINGON"`))
	require.EqualError(t, err, "invalid TransliterationMethodCode alias")
	// And a default value of "OTHR" is assigned to the object
	require.Equal(t, unknown, ivms101.TransliterationMethodCode(0))

	// Test that incorrect json input (no quotes) triggers a helpful error
	var badCode ivms101.TransliterationMethodCode
	err = badCode.UnmarshalJSON([]byte("KORE"))
	require.EqualError(t, err, "could not parse TransliterationMethodCode from value")
	// And a default value of "OTHR" is assigned to the object
	require.Equal(t, badCode, ivms101.TransliterationMethodCode(0))
}
