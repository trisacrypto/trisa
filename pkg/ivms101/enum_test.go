package ivms101_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

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
