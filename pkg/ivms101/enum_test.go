package ivms101_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func TestJSONMarshal(t *testing.T) {
	t.Run("NaturalPersonNameTypeCode", func(t *testing.T) {
		// Test an array of expected ivms101 protocol buffer enums
		// to inspect whether we correctly marshal them to the correct
		// compatible ivms101 format
		tests := map[ivms101.NaturalPersonNameTypeCode][]byte{
			0: []byte(`"MISC"`),
			1: []byte(`"ALIA"`),
			2: []byte(`"BIRT"`),
			3: []byte(`"MAID"`),
			4: []byte(`"LEGL"`),
		}

		for code, expected := range tests {
			actual, err := json.Marshal(code)
			require.NoError(t, err, "could not marshal %q", code.String())
			require.Equal(t, expected, actual, "incorrect marshal for test case %q", code.String())
		}
	})

	t.Run("LegalPersonNameTypeCode", func(t *testing.T) {
		// Test an array of expected ivms101 protocol buffer enums
		// to inspect whether we correctly marshal them to the correct
		// compatible ivms101 format
		tests := map[ivms101.LegalPersonNameTypeCode][]byte{
			0: []byte(`"MISC"`),
			1: []byte(`"LEGL"`),
			2: []byte(`"SHRT"`),
			3: []byte(`"TRAD"`),
		}

		for code, expected := range tests {
			actual, err := json.Marshal(code)
			require.NoError(t, err, "could not marshal %q", code.String())
			require.Equal(t, expected, actual, "incorrect marshal for test case %q", code.String())
		}
	})

	t.Run("AddressTypeCode", func(t *testing.T) {
		// Test an array of expected ivms101 protocol buffer enums
		// to inspect whether we correctly marshal them to the correct
		// compatible ivms101 format
		tests := map[ivms101.AddressTypeCode][]byte{
			0: []byte(`"MISC"`),
			1: []byte(`"HOME"`),
			2: []byte(`"BIZZ"`),
			3: []byte(`"GEOG"`),
		}

		for code, expected := range tests {
			actual, err := json.Marshal(code)
			require.NoError(t, err, "could not marshal %q", code.String())
			require.Equal(t, expected, actual, "incorrect marshal for test case %q", code.String())
		}
	})

	t.Run("NationalIdentifierTypeCode", func(t *testing.T) {
		// Test an array of expected ivms101 protocol buffer enums
		// to inspect whether we correctly marshal them to the correct
		// compatible ivms101 format
		tests := map[ivms101.NationalIdentifierTypeCode][]byte{
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

		for code, expected := range tests {
			actual, err := json.Marshal(code)
			require.NoError(t, err, "could not marshal %q", code.String())
			require.Equal(t, expected, actual, "incorrect marshal for test case %q", code.String())
		}
	})

	t.Run("TransliterationMethodCode", func(t *testing.T) {
		// Test an array of expected ivms101 protocol buffer enums
		// to inspect whether we correctly marshal them to the correct
		// compatible ivms101 format
		tests := map[ivms101.TransliterationMethodCode][]byte{
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

		for code, expected := range tests {
			actual, err := json.Marshal(code)
			require.NoError(t, err, "could not marshal %q", code.String())
			require.Equal(t, expected, actual, "incorrect marshal for test case %q", code.String())
		}
	})
}

func TestJSONUnmarshal(t *testing.T) {
	t.Run("NaturalPersonNameTypeCode", func(t *testing.T) {
		valid := []struct {
			expected ivms101.NaturalPersonNameTypeCode
			message  []byte
		}{
			{ivms101.NaturalPersonMisc, []byte(`"MISC"`)},
			{ivms101.NaturalPersonAlias, []byte(`"ALIA"`)},
			{ivms101.NaturalPersonBirth, []byte(`"BIRT"`)},
			{ivms101.NaturalPersonMaiden, []byte(`"MAID"`)},
			{ivms101.NaturalPersonLegal, []byte(`"LEGL"`)},
			{ivms101.NaturalPersonMisc, []byte(`0`)},
			{ivms101.NaturalPersonAlias, []byte(`1`)},
			{ivms101.NaturalPersonBirth, []byte(`2`)},
			{ivms101.NaturalPersonMaiden, []byte(`3`)},
			{ivms101.NaturalPersonLegal, []byte(`4`)},
			{ivms101.NaturalPersonMisc, []byte(`"NATURAL_PERSON_NAME_TYPE_CODE_MISC"`)},
			{ivms101.NaturalPersonAlias, []byte(`"NATURAL_PERSON_NAME_TYPE_CODE_ALIA"`)},
			{ivms101.NaturalPersonBirth, []byte(`"NATURAL_PERSON_NAME_TYPE_CODE_BIRT"`)},
			{ivms101.NaturalPersonMaiden, []byte(`"NATURAL_PERSON_NAME_TYPE_CODE_MAID"`)},
			{ivms101.NaturalPersonLegal, []byte(`"NATURAL_PERSON_NAME_TYPE_CODE_LEGL"`)},
			{ivms101.NaturalPersonMisc, []byte(`"misc"`)},
			{ivms101.NaturalPersonAlias, []byte(`"alia"`)},
			{ivms101.NaturalPersonBirth, []byte(`"birt"`)},
			{ivms101.NaturalPersonMaiden, []byte(`"maid"`)},
			{ivms101.NaturalPersonLegal, []byte(`"legl"`)},
		}

		for i, tc := range valid {
			var code ivms101.NaturalPersonNameTypeCode
			require.NoError(t, json.Unmarshal(tc.message, &code), "could not unmarshal data in test case %d", i)
			require.Equal(t, tc.expected, code, "incorrect unmarshal for test case %d", i)
		}

		invalid := []struct {
			target  error
			message []byte
		}{
			{ivms101.ErrInvalidNaturalPersonNameTypeCode, []byte(`"STAGE"`)},
			{ivms101.ErrInvalidNaturalPersonNameTypeCode, []byte("{}")},
		}

		for i, tc := range invalid {
			var code ivms101.NaturalPersonNameTypeCode
			err := json.Unmarshal(tc.message, &code)
			require.ErrorIs(t, err, tc.target, "expected error on test case %d", i)
			require.Zero(t, code, "expected zero valued code on test case %d", i)
		}
	})

	t.Run("LegalPersonNameTypeCode", func(t *testing.T) {
		valid := []struct {
			expected ivms101.LegalPersonNameTypeCode
			message  []byte
		}{
			{ivms101.LegalPersonMisc, []byte(`"MISC"`)},
			{ivms101.LegalPersonLegal, []byte(`"LEGL"`)},
			{ivms101.LegalPersonShort, []byte(`"SHRT"`)},
			{ivms101.LegalPersonTrading, []byte(`"TRAD"`)},
			{ivms101.LegalPersonMisc, []byte(`0`)},
			{ivms101.LegalPersonLegal, []byte(`1`)},
			{ivms101.LegalPersonShort, []byte(`2`)},
			{ivms101.LegalPersonTrading, []byte(`3`)},
			{ivms101.LegalPersonMisc, []byte(`"LEGAL_PERSON_NAME_TYPE_CODE_MISC"`)},
			{ivms101.LegalPersonLegal, []byte(`"LEGAL_PERSON_NAME_TYPE_CODE_LEGL"`)},
			{ivms101.LegalPersonShort, []byte(`"LEGAL_PERSON_NAME_TYPE_CODE_SHRT"`)},
			{ivms101.LegalPersonTrading, []byte(`"LEGAL_PERSON_NAME_TYPE_CODE_TRAD"`)},
			{ivms101.LegalPersonMisc, []byte(`"misc"`)},
			{ivms101.LegalPersonLegal, []byte(`"legl"`)},
			{ivms101.LegalPersonShort, []byte(`"shrt"`)},
			{ivms101.LegalPersonTrading, []byte(`"trad"`)},
		}

		for i, tc := range valid {
			var code ivms101.LegalPersonNameTypeCode
			require.NoError(t, json.Unmarshal(tc.message, &code), "could not unmarshal data in test case %d", i)
			require.Equal(t, tc.expected, code, "incorrect unmarshal for test case %d", i)
		}

		invalid := []struct {
			target  error
			message []byte
		}{
			{ivms101.ErrInvalidLegalPersonNameTypeCode, []byte(`"SHEL"`)},
			{ivms101.ErrInvalidLegalPersonNameTypeCode, []byte(`{}`)},
		}

		for i, tc := range invalid {
			var code ivms101.LegalPersonNameTypeCode
			err := json.Unmarshal(tc.message, &code)
			require.ErrorIs(t, err, tc.target, "expected error on test case %d", i)
			require.Zero(t, code, "expected zero valued code on test case %d", i)
		}
	})

	t.Run("AddressTypeCode", func(t *testing.T) {
		valid := []struct {
			expected ivms101.AddressTypeCode
			message  []byte
		}{
			{ivms101.AddressTypeMisc, []byte(`"MISC"`)},
			{ivms101.AddressTypeHome, []byte(`"HOME"`)},
			{ivms101.AddressTypeBusiness, []byte(`"BIZZ"`)},
			{ivms101.AddressTypeGeographic, []byte(`"GEOG"`)},
			{ivms101.AddressTypeMisc, []byte(`0`)},
			{ivms101.AddressTypeHome, []byte(`1`)},
			{ivms101.AddressTypeBusiness, []byte(`2`)},
			{ivms101.AddressTypeGeographic, []byte(`3`)},
			{ivms101.AddressTypeMisc, []byte(`"ADDRESS_TYPE_CODE_MISC"`)},
			{ivms101.AddressTypeHome, []byte(`"ADDRESS_TYPE_CODE_HOME"`)},
			{ivms101.AddressTypeBusiness, []byte(`"ADDRESS_TYPE_CODE_BIZZ"`)},
			{ivms101.AddressTypeGeographic, []byte(`"ADDRESS_TYPE_CODE_GEOG"`)},
			{ivms101.AddressTypeMisc, []byte(`"misc"`)},
			{ivms101.AddressTypeHome, []byte(`"home"`)},
			{ivms101.AddressTypeBusiness, []byte(`"bizz"`)},
			{ivms101.AddressTypeGeographic, []byte(`"geog"`)},
		}

		for i, tc := range valid {
			var code ivms101.AddressTypeCode
			require.NoError(t, json.Unmarshal(tc.message, &code), "could not unmarshal data in test case %d", i)
			require.Equal(t, tc.expected, code, "incorrect unmarshal for test case %d", i)
		}

		invalid := []struct {
			target  error
			message []byte
		}{
			{ivms101.ErrInvalidAddressTypeCode, []byte(`"LALA"`)},
			{ivms101.ErrInvalidAddressTypeCode, []byte(`{}`)},
		}

		for i, tc := range invalid {
			var code ivms101.AddressTypeCode
			err := json.Unmarshal(tc.message, &code)
			require.ErrorIs(t, err, tc.target, "expected error on test case %d", i)
			require.Zero(t, code, "expected zero valued code on test case %d", i)
		}
	})

	t.Run("NationalIdentifierTypeCode", func(t *testing.T) {
		valid := []struct {
			expected ivms101.NationalIdentifierTypeCode
			message  []byte
		}{
			{ivms101.NationalIdentifierMISC, []byte(`"MISC"`)},
			{ivms101.NationalIdentifierARNU, []byte(`"ARNU"`)},
			{ivms101.NationalIdentifierCCPT, []byte(`"CCPT"`)},
			{ivms101.NationalIdentifierRAID, []byte(`"RAID"`)},
			{ivms101.NationalIdentifierDRLC, []byte(`"DRLC"`)},
			{ivms101.NationalIdentifierFIIN, []byte(`"FIIN"`)},
			{ivms101.NationalIdentifierTXID, []byte(`"TXID"`)},
			{ivms101.NationalIdentifierSOCS, []byte(`"SOCS"`)},
			{ivms101.NationalIdentifierIDCD, []byte(`"IDCD"`)},
			{ivms101.NationalIdentifierLEIX, []byte(`"LEIX"`)},
			{ivms101.NationalIdentifierMISC, []byte(`0`)},
			{ivms101.NationalIdentifierARNU, []byte(`1`)},
			{ivms101.NationalIdentifierCCPT, []byte(`2`)},
			{ivms101.NationalIdentifierRAID, []byte(`3`)},
			{ivms101.NationalIdentifierDRLC, []byte(`4`)},
			{ivms101.NationalIdentifierFIIN, []byte(`5`)},
			{ivms101.NationalIdentifierTXID, []byte(`6`)},
			{ivms101.NationalIdentifierSOCS, []byte(`7`)},
			{ivms101.NationalIdentifierIDCD, []byte(`8`)},
			{ivms101.NationalIdentifierLEIX, []byte(`9`)},
			{ivms101.NationalIdentifierMISC, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_MISC"`)},
			{ivms101.NationalIdentifierARNU, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_ARNU"`)},
			{ivms101.NationalIdentifierCCPT, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_CCPT"`)},
			{ivms101.NationalIdentifierRAID, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_RAID"`)},
			{ivms101.NationalIdentifierDRLC, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_DRLC"`)},
			{ivms101.NationalIdentifierFIIN, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_FIIN"`)},
			{ivms101.NationalIdentifierTXID, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_TXID"`)},
			{ivms101.NationalIdentifierSOCS, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_SOCS"`)},
			{ivms101.NationalIdentifierIDCD, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_IDCD"`)},
			{ivms101.NationalIdentifierLEIX, []byte(`"NATIONAL_IDENTIFIER_TYPE_CODE_LEIX"`)},
			{ivms101.NationalIdentifierMISC, []byte(`"misc"`)},
			{ivms101.NationalIdentifierARNU, []byte(`"arnu"`)},
			{ivms101.NationalIdentifierCCPT, []byte(`"ccpt"`)},
			{ivms101.NationalIdentifierRAID, []byte(`"raid"`)},
			{ivms101.NationalIdentifierDRLC, []byte(`"drlc"`)},
			{ivms101.NationalIdentifierFIIN, []byte(`"fiin"`)},
			{ivms101.NationalIdentifierTXID, []byte(`"txid"`)},
			{ivms101.NationalIdentifierSOCS, []byte(`"socs"`)},
			{ivms101.NationalIdentifierIDCD, []byte(`"idcd"`)},
			{ivms101.NationalIdentifierLEIX, []byte(`"leix"`)},
		}

		for i, tc := range valid {
			var code ivms101.NationalIdentifierTypeCode
			require.NoError(t, json.Unmarshal(tc.message, &code), "could not unmarshal data in test case %d", i)
			require.Equal(t, tc.expected, code, "incorrect unmarshal for test case %d", i)
		}

		invalid := []struct {
			target  error
			message []byte
		}{
			{ivms101.ErrInvalidNationalIdentifierTypeCode, []byte(`"ACME"`)},
			{ivms101.ErrInvalidNationalIdentifierTypeCode, []byte(`{}`)},
		}

		for i, tc := range invalid {
			var code ivms101.NationalIdentifierTypeCode
			err := json.Unmarshal(tc.message, &code)
			require.ErrorIs(t, err, tc.target, "expected error on test case %d", i)
			require.Zero(t, code, "expected zero valued code on test case %d", i)
		}
	})

	t.Run("TransliterationMethodCode", func(t *testing.T) {
		valid := []struct {
			expected ivms101.TransliterationMethodCode
			message  []byte
		}{
			{ivms101.TransliterationMethodOTHR, []byte(`"OTHR"`)},
			{ivms101.TransliterationMethodARAB, []byte(`"ARAB"`)},
			{ivms101.TransliterationMethodARAN, []byte(`"ARAN"`)},
			{ivms101.TransliterationMethodARMN, []byte(`"ARMN"`)},
			{ivms101.TransliterationMethodCYRL, []byte(`"CYRL"`)},
			{ivms101.TransliterationMethodDEVA, []byte(`"DEVA"`)},
			{ivms101.TransliterationMethodGEOR, []byte(`"GEOR"`)},
			{ivms101.TransliterationMethodGREK, []byte(`"GREK"`)},
			{ivms101.TransliterationMethodHANI, []byte(`"HANI"`)},
			{ivms101.TransliterationMethodHEBR, []byte(`"HEBR"`)},
			{ivms101.TransliterationMethodKANA, []byte(`"KANA"`)},
			{ivms101.TransliterationMethodKORE, []byte(`"KORE"`)},
			{ivms101.TransliterationMethodTHAI, []byte(`"THAI"`)},
			{ivms101.TransliterationMethodOTHR, []byte(`0`)},
			{ivms101.TransliterationMethodARAB, []byte(`1`)},
			{ivms101.TransliterationMethodARAN, []byte(`2`)},
			{ivms101.TransliterationMethodARMN, []byte(`3`)},
			{ivms101.TransliterationMethodCYRL, []byte(`4`)},
			{ivms101.TransliterationMethodDEVA, []byte(`5`)},
			{ivms101.TransliterationMethodGEOR, []byte(`6`)},
			{ivms101.TransliterationMethodGREK, []byte(`7`)},
			{ivms101.TransliterationMethodHANI, []byte(`8`)},
			{ivms101.TransliterationMethodHEBR, []byte(`9`)},
			{ivms101.TransliterationMethodKANA, []byte(`10`)},
			{ivms101.TransliterationMethodKORE, []byte(`11`)},
			{ivms101.TransliterationMethodTHAI, []byte(`12`)},
			{ivms101.TransliterationMethodOTHR, []byte(`"TRANSLITERATION_METHOD_CODE_OTHR"`)},
			{ivms101.TransliterationMethodARAB, []byte(`"TRANSLITERATION_METHOD_CODE_ARAB"`)},
			{ivms101.TransliterationMethodARAN, []byte(`"TRANSLITERATION_METHOD_CODE_ARAN"`)},
			{ivms101.TransliterationMethodARMN, []byte(`"TRANSLITERATION_METHOD_CODE_ARMN"`)},
			{ivms101.TransliterationMethodCYRL, []byte(`"TRANSLITERATION_METHOD_CODE_CYRL"`)},
			{ivms101.TransliterationMethodDEVA, []byte(`"TRANSLITERATION_METHOD_CODE_DEVA"`)},
			{ivms101.TransliterationMethodGEOR, []byte(`"TRANSLITERATION_METHOD_CODE_GEOR"`)},
			{ivms101.TransliterationMethodGREK, []byte(`"TRANSLITERATION_METHOD_CODE_GREK"`)},
			{ivms101.TransliterationMethodHANI, []byte(`"TRANSLITERATION_METHOD_CODE_HANI"`)},
			{ivms101.TransliterationMethodHEBR, []byte(`"TRANSLITERATION_METHOD_CODE_HEBR"`)},
			{ivms101.TransliterationMethodKANA, []byte(`"TRANSLITERATION_METHOD_CODE_KANA"`)},
			{ivms101.TransliterationMethodKORE, []byte(`"TRANSLITERATION_METHOD_CODE_KORE"`)},
			{ivms101.TransliterationMethodTHAI, []byte(`"TRANSLITERATION_METHOD_CODE_THAI"`)},
			{ivms101.TransliterationMethodOTHR, []byte(`"othr"`)},
			{ivms101.TransliterationMethodARAB, []byte(`"arab"`)},
			{ivms101.TransliterationMethodARAN, []byte(`"aran"`)},
			{ivms101.TransliterationMethodARMN, []byte(`"armn"`)},
			{ivms101.TransliterationMethodCYRL, []byte(`"cyrl"`)},
			{ivms101.TransliterationMethodDEVA, []byte(`"deva"`)},
			{ivms101.TransliterationMethodGEOR, []byte(`"geor"`)},
			{ivms101.TransliterationMethodGREK, []byte(`"grek"`)},
			{ivms101.TransliterationMethodHANI, []byte(`"hani"`)},
			{ivms101.TransliterationMethodHEBR, []byte(`"hebr"`)},
			{ivms101.TransliterationMethodKANA, []byte(`"kana"`)},
			{ivms101.TransliterationMethodKORE, []byte(`"kore"`)},
			{ivms101.TransliterationMethodTHAI, []byte(`"thai"`)},
		}

		for i, tc := range valid {
			var code ivms101.TransliterationMethodCode
			require.NoError(t, json.Unmarshal(tc.message, &code), "could not unmarshal data in test case %d", i)
			require.Equal(t, tc.expected, code, "incorrect unmarshal for test case %d")
		}

		invalid := []struct {
			target  error
			message []byte
		}{
			{ivms101.ErrInvalidTransliterationMethodCode, []byte(`"KLINGON"`)},
			{ivms101.ErrInvalidTransliterationMethodCode, []byte(`{}`)},
		}

		for i, tc := range invalid {
			var code ivms101.TransliterationMethodCode
			err := json.Unmarshal(tc.message, &code)
			require.ErrorIs(t, err, tc.target, "expected error on test case %d", i)
			require.Zero(t, code, "expected zero valued code on test case %d", i)
		}
	})
}

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
			{"alias", ivms101.ErrInvalidNaturalPersonNameTypeCode},
			{2, ivms101.ErrInvalidNaturalPersonNameTypeCode},
			{"", ivms101.ErrInvalidNaturalPersonNameTypeCode},
			{"NATURAL_PERSON_NAME_TYPE_CODE_FOO", ivms101.ErrInvalidNaturalPersonNameTypeCode},
			{nil, ivms101.ErrInvalidNaturalPersonNameTypeCode},
			{int32(28), ivms101.ErrInvalidNaturalPersonNameTypeCode},
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
			{"legal", ivms101.ErrInvalidLegalPersonNameTypeCode},
			{2, ivms101.ErrInvalidLegalPersonNameTypeCode},
			{"", ivms101.ErrInvalidLegalPersonNameTypeCode},
			{"LEGAL_PERSON_NAME_TYPE_CODE_FOO", ivms101.ErrInvalidLegalPersonNameTypeCode},
			{nil, ivms101.ErrInvalidLegalPersonNameTypeCode},
			{int32(28), ivms101.ErrInvalidLegalPersonNameTypeCode},
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
			{"business", ivms101.ErrInvalidAddressTypeCode},
			{2, ivms101.ErrInvalidAddressTypeCode},
			{"", ivms101.ErrInvalidAddressTypeCode},
			{"ADDRESS_TYPE_CODE_FOO", ivms101.ErrInvalidAddressTypeCode},
			{nil, ivms101.ErrInvalidAddressTypeCode},
			{int32(28), ivms101.ErrInvalidAddressTypeCode},
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
			{"passport", ivms101.ErrInvalidNationalIdentifierTypeCode},
			{2, ivms101.ErrInvalidNationalIdentifierTypeCode},
			{"", ivms101.ErrInvalidNationalIdentifierTypeCode},
			{"NATIONAL_IDENTIFIER_TYPE_CODE_FOO", ivms101.ErrInvalidNationalIdentifierTypeCode},
			{nil, ivms101.ErrInvalidNationalIdentifierTypeCode},
			{int32(28), ivms101.ErrInvalidNationalIdentifierTypeCode},
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
			{"arabic", ivms101.ErrInvalidTransliterationMethodCode},
			{2, ivms101.ErrInvalidTransliterationMethodCode},
			{"", ivms101.ErrInvalidTransliterationMethodCode},
			{"TRANSLITERATION_METHOD_CODE_FOO", ivms101.ErrInvalidTransliterationMethodCode},
			{nil, ivms101.ErrInvalidTransliterationMethodCode},
			{int32(28), ivms101.ErrInvalidTransliterationMethodCode},
		}

		for i, tc := range invalidTestCases {
			actual, err := ivms101.ParseTransliterationMethodCode(tc.input)
			require.ErrorIs(t, err, tc.err, "expected error on invalid test case %d", i)
			require.Zero(t, actual, "expected zero value for invalid test case %d", i)
		}
	})
}
