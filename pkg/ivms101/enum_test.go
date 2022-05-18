package ivms101_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/nsf/jsondiff"
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
	tErr := errors.New("invalid NaturalPersonNameTypeCode alias")
	require.Equal(t, err, tErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.NaturalPersonNameTypeCode(0))

	// Test that incorrect json input triggers a helpful error
	var badCode ivms101.NaturalPersonNameTypeCode
	err = badCode.UnmarshalJSON([]byte("ALIA"))
	jErr := errors.New("could not parse NaturalPersonNameTypeCode from value")
	require.Equal(t, err, jErr)
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
	tErr := errors.New("invalid LegalPersonNameTypeCode alias")
	require.Equal(t, err, tErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.LegalPersonNameTypeCode(0))

	// Test that incorrect json input triggers a helpful error
	var badCode ivms101.LegalPersonNameTypeCode
	err = badCode.UnmarshalJSON([]byte("LEGL"))
	jErr := errors.New("could not parse LegalPersonNameTypeCode from value")
	require.Equal(t, err, jErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.LegalPersonNameTypeCode(0))
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
	tErr := errors.New("invalid NationalIdentifierTypeCode alias")
	require.Equal(t, err, tErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.NationalIdentifierTypeCode(0))

	// Test that incorrect json input triggers a helpful error
	var badCode ivms101.NationalIdentifierTypeCode
	err = badCode.UnmarshalJSON([]byte("LEIX"))
	jErr := errors.New("could not parse NationalIdentifierTypeCode from value")
	require.Equal(t, err, jErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.NationalIdentifierTypeCode(0))
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
	tErr := errors.New("invalid AddressTypeCode alias")
	require.Equal(t, err, tErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, unknown, ivms101.AddressTypeCode(0))

	// Test that incorrect json input triggers a helpful error
	var badCode ivms101.AddressTypeCode
	err = badCode.UnmarshalJSON([]byte("HOME"))
	jErr := errors.New("could not parse AddressTypeCode from value")
	require.Equal(t, err, jErr)
	// And a default value of "MISC" is assigned to the object
	require.Equal(t, badCode, ivms101.AddressTypeCode(0))
}

func TestMarshalNatNameID(t *testing.T) {
	// Test ivms101 protocol buffer struct to inspect whether we
	// correctly marshal to the correct compatible ivms101 format
	ident := &ivms101.NaturalPersonNameId{
		PrimaryIdentifier:   "clark",
		SecondaryIdentifier: "kent",
		NameIdentifierType:  ivms101.NaturalPersonNameTypeCode(1),
	}
	expected := []byte(`{"primaryIdentifier":"clark","secondaryIdentifier":"kent","nameIdentifierType":"ALIA"}`)
	compat, err := json.Marshal(ident)
	require.Nil(t, err)
	require.Equal(t, expected, compat)
}

func TestUnmarshalNatNameID(t *testing.T) {
	// Test compatible ivms101 identifiers to inspect whether we correctly
	// unmarshal them to the correct protocol buffer structs
	expected := &ivms101.NaturalPersonNameId{
		PrimaryIdentifier:   "clark",
		SecondaryIdentifier: "kent",
		NameIdentifierType:  ivms101.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_ALIA,
	}
	ident := &ivms101.NaturalPersonNameId{}
	compat := []byte(`{"nameIdentifierType":"ALIA","primaryIdentifier":"clark","secondaryIdentifier":"kent"}`)
	err := json.Unmarshal(compat, ident)
	require.Nil(t, err)
	require.Equal(t, expected, ident)
}

func TestMarshalLocNatNameID(t *testing.T) {
	// Test ivms101 protocol buffer struct to inspect whether we
	// correctly marshal to the correct compatible ivms101 format
	ident := &ivms101.LocalNaturalPersonNameId{
		PrimaryIdentifier:   "kal",
		SecondaryIdentifier: "el",
		NameIdentifierType:  ivms101.NaturalPersonNameTypeCode(2),
	}
	expected := []byte(`{"primaryIdentifier":"kal","secondaryIdentifier":"el","nameIdentifierType":"BIRT"}`)
	compat, err := json.Marshal(ident)
	require.Nil(t, err)
	require.Equal(t, expected, compat)
}

func TestUnmarshalLocNatNameID(t *testing.T) {
	// Test compatible ivms101 identifiers to inspect whether we correctly
	// unmarshal them to the correct protocol buffer structs
	expected := &ivms101.LocalNaturalPersonNameId{
		PrimaryIdentifier:   "kal",
		SecondaryIdentifier: "el",
		NameIdentifierType:  ivms101.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_BIRT,
	}
	ident := &ivms101.LocalNaturalPersonNameId{}
	compat := []byte(`{"primaryIdentifier":"kal","secondaryIdentifier":"el","nameIdentifierType":"BIRT"}`)
	err := json.Unmarshal(compat, ident)
	require.Nil(t, err)
	require.Equal(t, expected, ident)
}

func TestMarshalLegNameID(t *testing.T) {
	// Test ivms101 protocol buffer struct to inspect whether we
	// correctly marshal to the correct compatible ivms101 format
	ident := &ivms101.LegalPersonNameId{
		LegalPersonName:               "acme labs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode(1),
	}
	expected := []byte(`{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}`)
	compat, err := json.Marshal(ident)
	require.Nil(t, err)
	require.Equal(t, expected, compat)
}

func TestUnmarshalLegNameID(t *testing.T) {
	// Test compatible ivms101 identifiers to inspect whether we correctly
	// unmarshal them to the correct protocol buffer structs
	expected := &ivms101.LegalPersonNameId{
		LegalPersonName:               "acme labs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL,
	}
	ident := &ivms101.LegalPersonNameId{}
	compat := []byte(`{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}`)
	err := json.Unmarshal(compat, ident)
	require.Nil(t, err)
	require.Equal(t, expected, ident)
}

func TestMarshalNatName(t *testing.T) {
	// Test ivms101 protocol buffer struct to inspect whether we
	// correctly marshal to the correct compatible ivms101 format
	name := &ivms101.NaturalPersonNameId{
		PrimaryIdentifier:  "superman",
		NameIdentifierType: ivms101.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_LEGL,
	}
	birthname := &ivms101.LocalNaturalPersonNameId{
		PrimaryIdentifier:   "kal",
		SecondaryIdentifier: "el",
		NameIdentifierType:  ivms101.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_BIRT,
	}
	n := &ivms101.NaturalPersonName{
		NameIdentifiers:      []*ivms101.NaturalPersonNameId{name},
		LocalNameIdentifiers: []*ivms101.LocalNaturalPersonNameId{birthname},
	}
	expected := []byte(`{"localNameIdentifier":[{"primaryIdentifier":"kal","secondaryIdentifier":"el","nameIdentifierType":"BIRT"}],"nameIdentifier":[{"primaryIdentifier":"superman","nameIdentifierType":"LEGL"}]}`)
	compat, err := json.Marshal(n)
	require.Nil(t, err)
	require.Equal(t, expected, compat)
}

func TestUnmarshalNatName(t *testing.T) {
	// Test compatible ivms101 identifiers to inspect whether we correctly
	// unmarshal them to the correct protocol buffer structs
	name := &ivms101.NaturalPersonNameId{
		PrimaryIdentifier:  "superman",
		NameIdentifierType: ivms101.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_LEGL,
	}
	birthname := &ivms101.LocalNaturalPersonNameId{
		PrimaryIdentifier:   "kal",
		SecondaryIdentifier: "el",
		NameIdentifierType:  ivms101.NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_BIRT,
	}
	expected := &ivms101.NaturalPersonName{
		NameIdentifiers:      []*ivms101.NaturalPersonNameId{name},
		LocalNameIdentifiers: []*ivms101.LocalNaturalPersonNameId{birthname},
	}
	n := &ivms101.NaturalPersonName{}
	compat := []byte(`{"localNameIdentifier":[{"nameIdentifierType":"BIRT","primaryIdentifier":"kal","secondaryIdentifier":"el"}],"nameIdentifier":[{"nameIdentifierType":"LEGL","primaryIdentifier":"superman","secondaryIdentifier":""}],"phoneticNameIdentifier":null}`)
	err := json.Unmarshal(compat, n)
	require.Nil(t, err)
	require.Equal(t, expected, n)
}

func TestMarshalLegName(t *testing.T) {
	// Test ivms101 protocol buffer struct to inspect whether we
	// correctly marshal to the correct compatible ivms101 format
	name := &ivms101.LegalPersonNameId{
		LegalPersonName:               "acme labs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL,
	}
	tradename := &ivms101.LocalLegalPersonNameId{
		LegalPersonName:               "animaniacs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_TRAD,
	}
	n := &ivms101.LegalPersonName{
		NameIdentifiers:      []*ivms101.LegalPersonNameId{name},
		LocalNameIdentifiers: []*ivms101.LocalLegalPersonNameId{tradename},
	}
	expected := []byte(`{"localNameIdentifier":[{"legalPersonName":"animaniacs","legalPersonNameIdentifierType":"TRAD"}],"nameIdentifier":[{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}]}`)
	compat, err := json.Marshal(n)
	require.Nil(t, err)
	require.Equal(t, expected, compat)
}

func TestUnmarshalLegName(t *testing.T) {
	// Test compatible ivms101 identifiers to inspect whether we correctly
	// unmarshal them to the correct protocol buffer structs
	name := &ivms101.LegalPersonNameId{
		LegalPersonName:               "acme labs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL,
	}
	tradename := &ivms101.LocalLegalPersonNameId{
		LegalPersonName:               "animaniacs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_TRAD,
	}
	expected := &ivms101.LegalPersonName{
		NameIdentifiers:      []*ivms101.LegalPersonNameId{name},
		LocalNameIdentifiers: []*ivms101.LocalLegalPersonNameId{tradename},
	}
	n := &ivms101.LegalPersonName{}
	compat := []byte(`{"localNameIdentifier":[{"legalPersonName":"animaniacs","legalPersonNameIdentifierType":"TRAD"}],"nameIdentifier":[{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}],"phoneticNameIdentifier":null}`)
	err := json.Unmarshal(compat, n)
	require.Nil(t, err)
	require.Equal(t, expected, n)
}

func TestAddressMarshaling(t *testing.T) {
	// Should be able to unmarshal address with only address line
	addrLine := []byte(`{"addressLine":["4321 MmmmBop Lane","Middle America, USA","20000"]}`)
	var addrA *ivms101.Address
	require.NoError(t, json.Unmarshal(addrLine, &addrA))
	expected := &ivms101.Address{AddressLine: []string{"4321 MmmmBop Lane", "Middle America, USA", "20000"}}
	require.Equal(t, expected, addrA, "marshalled json differs from original")

	// Should be able to marshal address line back to PB struct
	compatA, err := json.Marshal(addrA)
	require.Nil(t, err)
	require.Equal(t, addrLine, compatA, "marshalled json differs from original")

	// Should be able to unmarshal a valid address from a JSON file
	data, err := ioutil.ReadFile("testdata/address.json")
	require.NoError(t, err)
	var addrB *ivms101.Address
	require.NoError(t, json.Unmarshal(data, &addrB))

	// Should be able to marshal to json without error
	compatB, err := json.Marshal(addrB)
	require.Nil(t, err)

	// Newly marshaled json should match original json
	diffOpts := jsondiff.DefaultConsoleOptions()
	res, _ := jsondiff.Compare(data, compatB, &diffOpts)
	require.Equal(t, res, jsondiff.FullMatch, "marshalled json differs from original")

}

// TODO
// func TestNaturalPersonMarshaling(t *testing.T) {}

// TODO
// func TestLegalPersonMarshaling(t *testing.T) {}
