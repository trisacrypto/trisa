package ivms101_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func TestLegalPerson(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/legalperson.json")
	require.NoError(t, err)

	// Should be able to load a valid legal person from JSON data
	var person *ivms101.LegalPerson
	require.NoError(t, json.Unmarshal(data, &person))

	// Legal person can't have a Country of Issue
	require.NoError(t, person.Validate())
	person.NationalIdentification.CountryOfIssue = "GB"
	require.Error(t, person.Validate())

	// If the National Identifier is an LEI, Registration Authority must be empty
	person.NationalIdentification.CountryOfIssue = ""
	person.NationalIdentification.NationalIdentifierType = 9
	person.NationalIdentification.RegistrationAuthority = ""
	require.NoError(t, person.Validate())
	person.NationalIdentification.RegistrationAuthority = "Chuck Norris"
	require.Error(t, person.Validate())

	// If the National Identifier is not an LEI, Registration Authority is mandatory
	person.NationalIdentification.NationalIdentifierType = 3
	require.NoError(t, person.Validate())
	person.NationalIdentification.RegistrationAuthority = ""
	require.Error(t, person.Validate())

	// Correct name can be validated
	require.NoError(t, person.Name.Validate())

	// Correctly structured name ids can be validated
	pid := &ivms101.LegalPersonNameId{
		LegalPersonName:               person.Name.NameIdentifiers[0].LegalPersonName,
		LegalPersonNameIdentifierType: person.Name.NameIdentifiers[0].LegalPersonNameIdentifierType,
	}
	require.NoError(t, pid.Validate())
	lpid := &ivms101.LocalLegalPersonNameId{
		LegalPersonName:               person.Name.NameIdentifiers[0].LegalPersonName,
		LegalPersonNameIdentifierType: person.Name.NameIdentifiers[0].LegalPersonNameIdentifierType,
	}
	require.NoError(t, lpid.Validate())

	// Complete address can be validated
	addr := &ivms101.Address{
		AddressType:    person.GeographicAddresses[0].AddressType,
		BuildingNumber: person.GeographicAddresses[0].BuildingNumber,
		StreetName:     person.GeographicAddresses[0].StreetName,
		TownName:       person.GeographicAddresses[0].TownName,
		PostCode:       person.GeographicAddresses[0].PostBox,
		Country:        person.GeographicAddresses[0].Country,
	}
	require.NoError(t, addr.Validate())

	// Correct national identifier can be validated
	require.NoError(t, person.NationalIdentification.Validate())

	// Should be able to convert a legal person into a generic Person
	gp := person.Person()
	require.Nil(t, gp.GetNaturalPerson())
	require.Equal(t, person, gp.GetLegalPerson())

	// Failure cases
	// Person with missing data should not produce a valid legal person
	notaperson := &ivms101.LegalPerson{}
	require.Error(t, notaperson.Validate())

	// Name with no identifiers can't be validated
	notaname := &ivms101.LegalPersonName{}
	require.Error(t, notaname.Validate())

	// Name missing type can't be validated
	noType := &ivms101.LegalPersonNameId{LegalPersonName: "Bob's Discount VASP, PLC"}
	notaname.NameIdentifiers = append(notaname.NameIdentifiers, noType)
	require.Error(t, notaname.Validate())

	// Incomplete name identifiers can't be validated
	pidBad := &ivms101.LegalPersonNameId{LegalPersonNameIdentifierType: 1}
	require.Error(t, pidBad.Validate())
	lpidBad := &ivms101.LocalLegalPersonNameId{LegalPersonNameIdentifierType: 1}
	require.Error(t, lpidBad.Validate())

	// Address with bad type can't be validated
	typeBad := &ivms101.Address{
		AddressType:    100000000,
		BuildingNumber: person.GeographicAddresses[0].BuildingNumber,
		StreetName:     person.GeographicAddresses[0].StreetName,
		TownName:       person.GeographicAddresses[0].TownName,
		PostCode:       person.GeographicAddresses[0].PostBox,
		Country:        person.GeographicAddresses[0].Country,
	}
	require.Error(t, typeBad.Validate())

	// Address with bad country can't be validated
	countryBad := &ivms101.Address{
		AddressType:    person.GeographicAddresses[0].AddressType,
		BuildingNumber: person.GeographicAddresses[0].BuildingNumber,
		StreetName:     person.GeographicAddresses[0].StreetName,
		TownName:       person.GeographicAddresses[0].TownName,
		PostCode:       person.GeographicAddresses[0].PostBox,
		Country:        "Lunar",
	}
	require.Error(t, countryBad.Validate())

	// Address with too many address lines can't be validated
	lines := []string{
		"123", "street", "lane", "road", "house", "cottage", "place", "usa",
	}
	linesBad := &ivms101.Address{
		AddressType: person.GeographicAddresses[0].AddressType,
		Country:     person.GeographicAddresses[0].Country,
		AddressLine: lines,
	}
	require.Error(t, linesBad.Validate())

	// Street is required if no address lines are provided
	noStreet := &ivms101.Address{
		AddressType:    person.GeographicAddresses[0].AddressType,
		BuildingNumber: person.GeographicAddresses[0].BuildingNumber,
		TownName:       person.GeographicAddresses[0].TownName,
		PostCode:       person.GeographicAddresses[0].PostBox,
		Country:        person.GeographicAddresses[0].Country,
	}
	require.Error(t, noStreet.Validate())

	// No national identification can't be validated
	noNid := &ivms101.NationalIdentification{NationalIdentifier: ""}
	require.Error(t, noNid.Validate())

	// Overlong national identification can't be validated
	longNid := &ivms101.NationalIdentification{
		NationalIdentifier:     "2343456987GHE97777342KIWERPM000000021287319636021864HT7450913054",
		NationalIdentifierType: 9,
		CountryOfIssue:         "GB",
		RegistrationAuthority:  "RA000589",
	}
	require.Error(t, longNid.Validate())

	// Invalid NID type can't be validated
	wrongNid := &ivms101.NationalIdentification{
		NationalIdentifier:     "213800AQUAUP6I215N33",
		NationalIdentifierType: 1000000000,
		CountryOfIssue:         "FR",
		RegistrationAuthority:  "RA000589",
	}
	require.Error(t, wrongNid.Validate())

	// Bad code for country of issue can't be validated
	badCode := &ivms101.NationalIdentification{
		NationalIdentifier:     "213800AQUAUP6I215N33",
		NationalIdentifierType: 4,
		CountryOfIssue:         "America",
		RegistrationAuthority:  "RA000589",
	}
	require.Error(t, badCode.Validate())
}

func TestNaturalPerson(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/naturalperson.json")
	require.NoError(t, err)

	// Should be able to load a valid natural person from JSON data
	var person *ivms101.NaturalPerson
	require.NoError(t, json.Unmarshal(data, &person))
	require.NoError(t, person.Validate())

	// Correct name can be validated
	require.NoError(t, person.Name.Validate())

	// Correct DOB can be validated
	require.NoError(t, person.DateAndPlaceOfBirth.Validate())

	// Correctly structured name ids can be validated
	pid := &ivms101.NaturalPersonNameId{
		PrimaryIdentifier:   person.Name.NameIdentifiers[0].PrimaryIdentifier,
		SecondaryIdentifier: person.Name.NameIdentifiers[0].SecondaryIdentifier,
		NameIdentifierType:  person.Name.NameIdentifiers[0].NameIdentifierType,
	}
	require.NoError(t, pid.Validate())
	lpid := &ivms101.LocalNaturalPersonNameId{
		PrimaryIdentifier:   person.Name.NameIdentifiers[0].PrimaryIdentifier,
		SecondaryIdentifier: person.Name.NameIdentifiers[0].SecondaryIdentifier,
		NameIdentifierType:  person.Name.NameIdentifiers[0].NameIdentifierType,
	}
	require.NoError(t, lpid.Validate())

	// Complete address can be validated
	addr := &ivms101.Address{
		AddressType:    person.GeographicAddresses[0].AddressType,
		BuildingNumber: person.GeographicAddresses[0].BuildingNumber,
		StreetName:     person.GeographicAddresses[0].StreetName,
		TownName:       person.GeographicAddresses[0].TownName,
		PostCode:       person.GeographicAddresses[0].PostBox,
		Country:        person.GeographicAddresses[0].Country,
	}
	require.NoError(t, addr.Validate())

	// Correct national identifier can be validated
	require.NoError(t, person.NationalIdentification.Validate())

	// Should be able to convert a legal person into a generic Person
	gp := person.Person()
	require.Nil(t, gp.GetLegalPerson())
	require.Equal(t, person, gp.GetNaturalPerson())

	// Failure cases
	// JSON data missing required fields should not produce a valid natural person
	notaperson := &ivms101.NaturalPerson{}
	require.Error(t, notaperson.Validate())

	// Name with no identifiers can't be validated
	notaname := &ivms101.NaturalPersonName{}
	require.Error(t, notaname.Validate())

	// Name missing type can't be validated
	noType := &ivms101.NaturalPersonNameId{
		PrimaryIdentifier:   "Annie",
		SecondaryIdentifier: "Oakley",
	}
	notaname.NameIdentifiers = append(notaname.NameIdentifiers, noType)
	require.Error(t, notaname.Validate())

	// Incomplete name identifiers can't be validated
	pidBad := &ivms101.NaturalPersonNameId{NameIdentifierType: 1}
	require.Error(t, pidBad.Validate())
	lpidBad := &ivms101.LocalNaturalPersonNameId{NameIdentifierType: 1}
	require.Error(t, lpidBad.Validate())

	// No date of birth can't be validated
	pobOnly := &ivms101.DateAndPlaceOfBirth{PlaceOfBirth: "Champaign, IL"}
	require.Error(t, pobOnly.Validate())

	// No place of birth can't be validated
	dobOnly := &ivms101.DateAndPlaceOfBirth{DateOfBirth: "2011-05-21"}
	require.Error(t, dobOnly.Validate())

	// Date must be parsable or can't be validated
	badDob := &ivms101.DateAndPlaceOfBirth{DateOfBirth: "80-80-80-80"}
	require.Error(t, badDob.Validate())

	// Can't be born in the future
	futureDob := &ivms101.DateAndPlaceOfBirth{DateOfBirth: "8000-05-21"}
	require.Error(t, futureDob.Validate())
}
