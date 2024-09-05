package ivms101_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

//
// Person JSON
//

func TestPersonMarshaling(t *testing.T) {
	data, err := os.ReadFile("testdata/person_natural_person.json")
	require.NoError(t, err, "could not load person with natural person json fixture")

	person := &ivms101.Person{}
	require.NoError(t, json.Unmarshal(data, person), "could not unmarshal person")
	require.Empty(t, person.GetLegalPerson(), "legal person returned from person with natural person fixture")
	require.NotEmpty(t, person.GetNaturalPerson(), "no natural person returned from person with natural person fixture")

	data, err = os.ReadFile("testdata/person_legal_person.json")
	require.NoError(t, err, "could not load person with legal person json fixture")

	person = &ivms101.Person{}
	require.NoError(t, json.Unmarshal(data, person), "could not unmarshal person")
	require.Empty(t, person.GetNaturalPerson(), "natural person returned from person with legal person fixture")
	require.NotEmpty(t, person.GetLegalPerson(), "no legal person returned from person with legal person fixture")

	data, err = os.ReadFile("testdata/person_both_persons.json")
	require.NoError(t, err, "could not load person with both persons json fixture")

	person = &ivms101.Person{}
	require.ErrorIs(t, json.Unmarshal(data, person), ivms101.ErrPersonOneOfViolation)
}

//
// NaturalPerson JSON
//

func TestNaturalPersonMarshaling(t *testing.T) {
	data, err := os.ReadFile("testdata/natural_person.json")
	require.NoError(t, err, "could not load natural person json fixture")

	var person *ivms101.NaturalPerson
	require.NoError(t, json.Unmarshal(data, &person), "could not unmarshal natural person")
	require.NotEmpty(t, person, "no lnatural person was unmarshaled")
	require.NotEmpty(t, person.Name, "no natural person name was unmarshaled")
	require.Len(t, person.Name.NameIdentifiers, 1, "incorrect number of name identifiers unmarshaled")
	require.NotEmpty(t, person.GeographicAddresses, "no natural person geographic addresses unmarshaled")
	require.Len(t, person.GeographicAddresses, 1, "incorrect number of geographic addresses unmarshaled")
	require.Equal(t, "1234abc", person.CustomerIdentification, "natural person customer identification not unmarshaled")
	require.NotEmpty(t, person.NationalIdentification, "no natural person national identification unmarshaled")
	require.Equal(t, "US", person.CountryOfResidence, "natural person country of residence not unmarshaled")

	compat, err := json.Marshal(person)
	require.NoError(t, err, "could not marshal natural person")

	require.JSONEq(t, string(data), string(compat), "marshalled json differs from original")
}

//
// NaturalPersonName and NaturalPersonNameIdentifiers JSON
//

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
	expected := `{"localNameIdentifier":[{"primaryIdentifier":"kal","secondaryIdentifier":"el","nameIdentifierType":"BIRT"}],"nameIdentifier":[{"primaryIdentifier":"superman","nameIdentifierType":"LEGL"}]}`
	compat, err := json.Marshal(n)
	require.Nil(t, err)
	require.JSONEq(t, expected, string(compat))
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

//
// Address JSON
//

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
	data, err := os.ReadFile("testdata/address.json")
	require.NoError(t, err)
	var addrB *ivms101.Address
	require.NoError(t, json.Unmarshal(data, &addrB))

	// Should be able to marshal to json without error
	compatB, err := json.Marshal(addrB)
	require.Nil(t, err)

	// Newly marshaled json should match original json
	require.JSONEq(t, string(data), string(compatB), "marshalled json differs from original")
}

//
// DateAndPlaceOfBirth JSON
//

func TestDateAndPlaceOfBirthMarshaling(t *testing.T) {
	dobData := []byte(`{"dateOfBirth": "1984-04-20", "placeOfBirth":"Montgomery, AL, USA"}`)
	var dob *ivms101.DateAndPlaceOfBirth
	require.NoError(t, json.Unmarshal(dobData, &dob))
	require.Equal(t, dob.DateOfBirth, "1984-04-20")
	require.Equal(t, dob.PlaceOfBirth, "Montgomery, AL, USA")

	outdata, err := json.Marshal(dob)
	require.NoError(t, err, "could not marshal date and place of birth")

	require.JSONEq(t, string(dobData), string(outdata), "marshalled json differs from original")
}

//
// NationalIdentification JSON
//

func TestNationalIdentificationMarshaling(t *testing.T) {
	data, err := os.ReadFile("testdata/national_identification.json")
	require.NoError(t, err, "could not load national identification json fixture")

	var natId *ivms101.NationalIdentification
	require.NoError(t, json.Unmarshal(data, &natId), "could not unmarshal national identification")
	require.Equal(t, "815026352", natId.NationalIdentifier)
	require.Equal(t, ivms101.NationalIdentifierDRLC, natId.NationalIdentifierType)
	require.Equal(t, "TV", natId.CountryOfIssue)
	require.Equal(t, "RA777777", natId.RegistrationAuthority)

	compat, err := json.Marshal(natId)
	require.NoError(t, err, "could not marshal national identification")

	require.JSONEq(t, string(data), string(compat), "marshalled json differs from original")
}

//
// LegalPerson JSON
//

func TestLegalPersonMarshaling(t *testing.T) {
	data, err := os.ReadFile("testdata/legal_person.json")
	require.NoError(t, err, "could not load legal person json fixture")

	var person *ivms101.LegalPerson
	require.NoError(t, json.Unmarshal(data, &person), "could not unmarshal legal person")
	require.NotEmpty(t, person, "no legal person was unmarshaled")
	require.NotEmpty(t, person.Name, "no legal person name was unmarshaled")
	require.Len(t, person.Name.NameIdentifiers, 2, "incorrect number of name identifiers unmarshaled")
	require.NotEmpty(t, person.GeographicAddresses, "no legal person geographic addresses unmarshaled")
	require.Len(t, person.GeographicAddresses, 1, "incorrect number of geographic addresses unmarshaled")
	require.Equal(t, "abc1234", person.CustomerNumber, "legal person customer number not unmarshaled")
	require.NotEmpty(t, person.NationalIdentification, "no legal person national identification unmarshaled")
	require.Equal(t, "GB", person.CountryOfRegistration, "legal person country of registration not unmarshaled")

	compat, err := json.Marshal(person)
	require.NoError(t, err, "could not marshal legal person")

	require.JSONEq(t, string(data), string(compat), "marshalled json differs from original")
}

//
// LegalPersonName and LegalPersonNameIdentifiers JSON
//

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
	expected := `{"localNameIdentifier":[{"legalPersonName":"animaniacs","legalPersonNameIdentifierType":"TRAD"}],"nameIdentifier":[{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}]}`
	compat, err := json.Marshal(n)
	require.Nil(t, err)
	require.JSONEq(t, expected, string(compat))
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

func TestMarshalLocalLegNameID(t *testing.T) {
	// Test ivms101 protocol buffer struct to inspect whether we
	// correctly marshal to the correct compatible ivms101 format
	ident := &ivms101.LocalLegalPersonNameId{
		LegalPersonName:               "acme labs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode(1),
	}
	expected := []byte(`{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}`)
	compat, err := json.Marshal(ident)
	require.Nil(t, err)
	require.Equal(t, expected, compat)
}

func TestUnmarshalLocalLegNameID(t *testing.T) {
	// Test compatible ivms101 identifiers to inspect whether we correctly
	// unmarshal them to the correct protocol buffer structs
	expected := &ivms101.LocalLegalPersonNameId{
		LegalPersonName:               "acme labs",
		LegalPersonNameIdentifierType: ivms101.LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL,
	}
	ident := &ivms101.LocalLegalPersonNameId{}
	compat := []byte(`{"legalPersonName":"acme labs","legalPersonNameIdentifierType":"LEGL"}`)
	err := json.Unmarshal(compat, ident)
	require.Nil(t, err)
	require.Equal(t, expected, ident)
}
