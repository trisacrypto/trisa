package ivms101

import (
	"encoding/json"
)

//===========================================================================
// Person Methods
//===========================================================================

type serialPerson struct {
	NaturalPerson *NaturalPerson `json:"naturalPerson,omitempty"`
	LegalPerson   *LegalPerson   `json:"legalPerson,omitempty"`
}

var serialPersonKeys = map[string]string{
	"naturalPerson":  "naturalPerson",
	"natural_person": "naturalPerson",
	"legalPerson":    "legalPerson",
	"legal_person":   "legalPerson",
}

func (p *Person) MarshalJSON() ([]byte, error) {
	middle := serialPerson{
		NaturalPerson: p.GetNaturalPerson(),
		LegalPerson:   p.GetLegalPerson(),
	}
	return json.Marshal(middle)
}

func (p *Person) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialPersonKeys); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := &serialPerson{}
	if err = json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Check oneof constraint
	if middle.NaturalPerson != nil && middle.LegalPerson != nil {
		return ErrPersonOneOfViolation
	}

	// Populate the person value
	switch {
	case middle.NaturalPerson != nil:
		p.Person = &Person_NaturalPerson{
			NaturalPerson: middle.NaturalPerson,
		}
	case middle.LegalPerson != nil:
		p.Person = &Person_LegalPerson{
			LegalPerson: middle.LegalPerson,
		}
	}
	return nil
}

//===========================================================================
// NaturalPerson Methods
//===========================================================================

// Person converts a NaturalPerson into a Person protobuf message type.
func (p *NaturalPerson) Person() *Person {
	return &Person{
		Person: &Person_NaturalPerson{
			NaturalPerson: p,
		},
	}
}

type serialNaturalPerson struct {
	Name               *NaturalPersonName      `json:"name,omitempty"`
	Address            []*Address              `json:"geographicAddress,omitempty"`
	Identification     *NationalIdentification `json:"nationalIdentification,omitempty"`
	CustomerID         string                  `json:"customerIdentification,omitempty"`
	DOB                *DateAndPlaceOfBirth    `json:"dateAndPlaceOfBirth,omitempty"`
	CountryOfResidence string                  `json:"countryOfResidence,omitempty"`
}

var serialNaturalPersonKeys = map[string]string{
	"name":                    "name",
	"names":                   "name",
	"geographicAddress":       "geographicAddress",
	"geographicAddresses":     "geographicAddress",
	"geographic_address":      "geographicAddress",
	"geographic_addresses":    "geographicAddress",
	"addresses":               "geographicAddress",
	"address":                 "geographicAddress",
	"nationalIdentification":  "nationalIdentification",
	"national_identification": "nationalIdentification",
	"customerIdentification":  "customerIdentification",
	"customer_identification": "customerIdentification",
	"dateAndPlaceOfBirth":     "dateAndPlaceOfBirth",
	"date_and_place_of_birth": "dateAndPlaceOfBirth",
	"dob":                     "dateAndPlaceOfBirth",
	"countryOfResidence":      "countryOfResidence",
	"country_of_residence":    "countryOfResidence",
	"country":                 "countryOfResidence",
}

func (n *NaturalPerson) MarshalJSON() ([]byte, error) {
	middle := serialNaturalPerson{
		Name:               n.Name,
		Address:            n.GeographicAddresses,
		Identification:     n.NationalIdentification,
		CustomerID:         n.CustomerIdentification,
		DOB:                n.DateAndPlaceOfBirth,
		CountryOfResidence: n.CountryOfResidence,
	}
	return json.Marshal(middle)
}

func (n *NaturalPerson) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialNaturalPersonKeys); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialNaturalPerson{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate the natural person value
	n.Name = middle.Name
	n.GeographicAddresses = middle.Address
	n.NationalIdentification = middle.Identification
	n.CustomerIdentification = middle.CustomerID
	n.DateAndPlaceOfBirth = middle.DOB
	n.CountryOfResidence = middle.CountryOfResidence

	return nil
}

//===========================================================================
// NaturalPersonNameIdentifiers Methods
//===========================================================================

type serialNaturalPersonName struct {
	NameIdentifiers         []*NaturalPersonNameId      `json:"nameIdentifier,omitempty"`
	LocalNameIdentifiers    []*LocalNaturalPersonNameId `json:"localNameIdentifier,omitempty"`
	PhoneticNameIdentifiers []*LocalNaturalPersonNameId `json:"phoneticNameIdentifier,omitempty"`
}

var serialNaturalPersonNameFields = map[string]string{
	"nameIdentifier":            "nameIdentifier",
	"nameIdentifiers":           "nameIdentifier",
	"name_identifier":           "nameIdentifier",
	"name_identifiers":          "nameIdentifier",
	"localNameIdentifier":       "localNameIdentifier",
	"localNameIdentifiers":      "localNameIdentifier",
	"local_name_identifier":     "localNameIdentifier",
	"local_name_identifiers":    "localNameIdentifier",
	"phoneticNameIdentifier":    "phoneticNameIdentifier",
	"phoneticNameIdentifiers":   "phoneticNameIdentifier",
	"phonetic_name_identifier":  "phoneticNameIdentifier",
	"phonetic_name_identifiers": "phoneticNameIdentifier",
}

func (n *NaturalPersonName) MarshalJSON() ([]byte, error) {
	middle := serialNaturalPersonName{
		LocalNameIdentifiers:    n.LocalNameIdentifiers,
		NameIdentifiers:         n.NameIdentifiers,
		PhoneticNameIdentifiers: n.PhoneticNameIdentifiers,
	}
	return json.Marshal(middle)
}

func (n *NaturalPersonName) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialNaturalPersonNameFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialNaturalPersonName{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate the natural person name value
	n.NameIdentifiers = middle.NameIdentifiers
	n.LocalNameIdentifiers = middle.LocalNameIdentifiers
	n.PhoneticNameIdentifiers = middle.PhoneticNameIdentifiers

	return nil
}

//===========================================================================
// NaturalPersonNameID Methods
//===========================================================================

type serialNaturalPersonNameID struct {
	PrimaryIdentifier   string                    `json:"primaryIdentifier,omitempty"`
	SecondaryIdentifier string                    `json:"secondaryIdentifier,omitempty"`
	NameIdentifierType  NaturalPersonNameTypeCode `json:"nameIdentifierType,omitempty"`
}

var serialNaturalPersonNameIDFields = map[string]string{
	"primaryIdentifier":    "primaryIdentifier",
	"primary_identifier":   "primaryIdentifier",
	"last_name":            "primaryIdentifier",
	"lastName":             "primaryIdentifier",
	"surname":              "primaryIdentifier",
	"family_name":          "primaryIdentifier",
	"familyName":           "primaryIdentifier",
	"secondaryIdentifier":  "secondaryIdentifier",
	"secondary_identifier": "secondaryIdentifier",
	"first_name":           "secondaryIdentifier",
	"firstName":            "secondaryIdentifier",
	"nameIdentifierType":   "nameIdentifierType",
	"name_identifier_type": "nameIdentifierType",
}

func (p *NaturalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialNaturalPersonNameID{
		PrimaryIdentifier:   p.PrimaryIdentifier,
		SecondaryIdentifier: p.SecondaryIdentifier,
		NameIdentifierType:  p.NameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *NaturalPersonNameId) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialNaturalPersonNameIDFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialNaturalPersonNameID{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate the natural person name id values
	p.PrimaryIdentifier = middle.PrimaryIdentifier
	p.SecondaryIdentifier = middle.SecondaryIdentifier
	p.NameIdentifierType = middle.NameIdentifierType

	return nil
}

//===========================================================================
// LocalNaturalPersonNameID Methods
//===========================================================================

type serialLocalNaturalPersonNameID struct {
	PrimaryIdentifier   string                    `json:"primaryIdentifier,omitempty"`
	SecondaryIdentifier string                    `json:"secondaryIdentifier,omitempty"`
	NameIdentifierType  NaturalPersonNameTypeCode `json:"nameIdentifierType,omitempty"`
}

var serialLocalNaturalPersonNameIDFields = map[string]string{
	"primaryIdentifier":    "primaryIdentifier",
	"primary_identifier":   "primaryIdentifier",
	"secondaryIdentifier":  "secondaryIdentifier",
	"secondary_identifier": "secondaryIdentifier",
	"nameIdentifierType":   "nameIdentifierType",
	"name_identifier_type": "nameIdentifierType",
}

func (p *LocalNaturalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialLocalNaturalPersonNameID{
		PrimaryIdentifier:   p.PrimaryIdentifier,
		SecondaryIdentifier: p.SecondaryIdentifier,
		NameIdentifierType:  p.NameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LocalNaturalPersonNameId) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialLocalNaturalPersonNameIDFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialLocalNaturalPersonNameID{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate the natural person name id values
	p.PrimaryIdentifier = middle.PrimaryIdentifier
	p.SecondaryIdentifier = middle.SecondaryIdentifier
	p.NameIdentifierType = middle.NameIdentifierType

	return nil
}

//===========================================================================
// Address Methods
//===========================================================================

type serialAddress struct {
	AddressType        AddressTypeCode `json:"addressType,omitempty"`
	Department         string          `json:"department,omitempty"`
	SubDepartment      string          `json:"subDepartment,omitempty"`
	StreetName         string          `json:"streetName,omitempty"`
	BuildingNumber     string          `json:"buildingNumber,omitempty"`
	BuildingName       string          `json:"buildingName,omitempty"`
	Floor              string          `json:"floor,omitempty"`
	PostBox            string          `json:"postBox,omitempty"`
	Room               string          `json:"room,omitempty"`
	PostCode           string          `json:"postCode,omitempty"`
	TownName           string          `json:"townName,omitempty"`
	TownLocationName   string          `json:"townLocationName,omitempty"`
	DistrictName       string          `json:"districtName,omitempty"`
	CountrySubDivision string          `json:"countrySubDivision,omitempty"`
	AddressLine        []string        `json:"addressLine,omitempty"`
	Country            string          `json:"country,omitempty"`
}

var serialAddressFields = map[string]string{
	"addressType":          "addressType",
	"address_type":         "addressType",
	"department":           "department",
	"subDepartment":        "subDepartment",
	"sub_department":       "subDepartment",
	"subdepartment":        "subDepartment",
	"streetName":           "streetName",
	"street_name":          "streetName",
	"street":               "streetName",
	"buildingNumber":       "buildingNumber",
	"building_number":      "buildingNumber",
	"number":               "buildingNumber",
	"buildingName":         "buildingName",
	"building_name":        "buildingName",
	"building":             "buildingName",
	"floor":                "floor",
	"postBox":              "postBox",
	"post_box":             "postBox",
	"pob":                  "postBox",
	"room":                 "room",
	"postCode":             "postCode",
	"post_code":            "postCode",
	"postalCode":           "postCode",
	"postal_code":          "postCode",
	"zipCode":              "postCode",
	"zip_code":             "postCode",
	"townName":             "townName",
	"town_name":            "townName",
	"town":                 "townName",
	"city":                 "townName",
	"townLocationName":     "townLocationName",
	"town_location_name":   "townLocationName",
	"locationName":         "townLocationName",
	"location_name":        "townLocationName",
	"districtName":         "districtName",
	"district_name":        "districtName",
	"district":             "districtName",
	"countrySubDivision":   "countrySubDivision",
	"country_sub_division": "countrySubDivision",
	"country_Subdivision":  "countrySubDivision",
	"country_subdivision":  "countrySubDivision",
	"state":                "countrySubDivision",
	"province":             "countrySubDivision",
	"addressLine":          "addressLine",
	"addressLines":         "addressLine",
	"address_line":         "addressLine",
	"address_lines":        "addressLine",
	"country":              "country",
}

func (a *Address) MarshalJSON() ([]byte, error) {
	middle := serialAddress{
		AddressType:        a.AddressType,
		Department:         a.Department,
		SubDepartment:      a.SubDepartment,
		StreetName:         a.StreetName,
		BuildingNumber:     a.BuildingNumber,
		BuildingName:       a.BuildingName,
		Floor:              a.Floor,
		PostBox:            a.PostBox,
		Room:               a.Room,
		PostCode:           a.PostCode,
		TownName:           a.TownName,
		TownLocationName:   a.TownLocationName,
		DistrictName:       a.DistrictName,
		CountrySubDivision: a.CountrySubDivision,
		AddressLine:        a.AddressLine,
		Country:            a.Country,
	}
	return json.Marshal(middle)
}

func (a *Address) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialAddressFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialAddress{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate address values
	a.AddressType = middle.AddressType
	a.Department = middle.Department
	a.SubDepartment = middle.SubDepartment
	a.StreetName = middle.StreetName
	a.BuildingNumber = middle.BuildingNumber
	a.BuildingName = middle.BuildingName
	a.Floor = middle.Floor
	a.PostBox = middle.PostBox
	a.Room = middle.Room
	a.PostCode = middle.PostCode
	a.TownName = middle.TownName
	a.TownLocationName = middle.TownLocationName
	a.DistrictName = middle.DistrictName
	a.CountrySubDivision = middle.CountrySubDivision
	a.AddressLine = middle.AddressLine
	a.Country = middle.Country

	return nil
}

//===========================================================================
// DateAndPlaceOfBirth Methods
//===========================================================================

type serialDateAndPlaceOfBirth struct {
	DateOfBirth  string `json:"dateOfBirth,omitempty"`
	PlaceOfBirth string `json:"placeOfBirth,omitempty"`
}

var serialDateAndPlaceOfBirthFields = map[string]string{
	"dateOfBirth":    "dateOfBirth",
	"date_of_birth":  "dateOfBirth",
	"dob":            "dateOfBirth",
	"placeOfBirth":   "placeOfBirth",
	"place_of_birth": "placeOfBirth",
	"pob":            "placeOfBirth",
}

func (d *DateAndPlaceOfBirth) MarshalJSON() ([]byte, error) {
	middle := serialDateAndPlaceOfBirth{
		DateOfBirth:  d.DateOfBirth,
		PlaceOfBirth: d.PlaceOfBirth,
	}
	return json.Marshal(middle)
}

func (d *DateAndPlaceOfBirth) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialDateAndPlaceOfBirthFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialDateAndPlaceOfBirth{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate date and place of birth values
	d.DateOfBirth = middle.DateOfBirth
	d.PlaceOfBirth = middle.PlaceOfBirth

	return nil
}

//===========================================================================
// NationalIdentification Methods
//===========================================================================

type serialNationalIdentification struct {
	NationalIdentifier     string                     `json:"nationalIdentifier,omitempty"`
	NationalIdentifierType NationalIdentifierTypeCode `json:"nationalIdentifierType,omitempty"`
	CountryOfIssue         string                     `json:"countryOfIssue,omitempty"`
	RegistrationAuthority  string                     `json:"registrationAuthority,omitempty"`
}

var serialNationalIdentificationFields = map[string]string{
	"nationalIdentifier":       "nationalIdentifier",
	"national_identifier":      "nationalIdentifier",
	"number":                   "nationalIdentifier",
	"identifierNumber":         "nationalIdentifier",
	"identifier_number":        "nationalIdentifier",
	"nationalIdentifierType":   "nationalIdentifierType",
	"national_identifier_type": "nationalIdentifierType",
	"countryOfIssue":           "countryOfIssue",
	"country_of_issue":         "countryOfIssue",
	"country":                  "countryOfIssue",
	"registrationAuthority":    "registrationAuthority",
	"registration_authority":   "registrationAuthority",
	"ra":                       "registrationAuthority",
}

func (n *NationalIdentification) MarshalJSON() ([]byte, error) {
	middle := serialNationalIdentification{
		NationalIdentifier:     n.NationalIdentifier,
		NationalIdentifierType: n.NationalIdentifierType,
		CountryOfIssue:         n.CountryOfIssue,
		RegistrationAuthority:  n.RegistrationAuthority,
	}
	return json.Marshal(middle)
}

func (n *NationalIdentification) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialNationalIdentificationFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialNationalIdentification{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate national identification values
	n.NationalIdentifier = middle.NationalIdentifier
	n.NationalIdentifierType = middle.NationalIdentifierType
	n.CountryOfIssue = middle.CountryOfIssue
	n.RegistrationAuthority = middle.RegistrationAuthority

	return nil
}

//===========================================================================
// LegalPerson Methods
//===========================================================================

type serialLegalPerson struct {
	Name                   *LegalPersonName        `json:"name,omitempty"`
	Address                []*Address              `json:"geographicAddress,omitempty"`
	CustomerNumber         string                  `json:"customerNumber,omitempty"`
	NationalIdentification *NationalIdentification `json:"nationalIdentification,omitempty"`
	CountryOfRegistration  string                  `json:"countryOfRegistration,omitempty"`
}

var serialLegalPersonFields = map[string]string{
	"name":                    "name",
	"names":                   "name",
	"geographicAddress":       "geographicAddress",
	"geographicAddresses":     "geographicAddress",
	"geographic_address":      "geographicAddress",
	"geographic_addresses":    "geographicAddress",
	"addresses":               "geographicAddress",
	"address":                 "geographicAddress",
	"customerNumber":          "customerNumber",
	"customer_number":         "customerNumber",
	"nationalIdentification":  "nationalIdentification",
	"national_identification": "nationalIdentification",
	"countryOfRegistration":   "countryOfRegistration",
	"country_of_registration": "countryOfRegistration",
	"country":                 "countryOfRegistration",
}

func (l *LegalPerson) MarshalJSON() ([]byte, error) {
	middle := serialLegalPerson{
		Name:                   l.Name,
		Address:                l.GeographicAddresses,
		CustomerNumber:         l.CustomerNumber,
		NationalIdentification: l.NationalIdentification,
		CountryOfRegistration:  l.CountryOfRegistration,
	}
	return json.Marshal(middle)
}

func (l *LegalPerson) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialLegalPersonFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialLegalPerson{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate legal person values
	l.Name = middle.Name
	l.GeographicAddresses = middle.Address
	l.CustomerNumber = middle.CustomerNumber
	l.NationalIdentification = middle.NationalIdentification
	l.CountryOfRegistration = middle.CountryOfRegistration

	return nil
}

//===========================================================================
// LegalPersonName Methods
//===========================================================================

// Person converts a LegalPerson into a Person protobuf message type.
func (p *LegalPerson) Person() *Person {
	return &Person{
		Person: &Person_LegalPerson{
			LegalPerson: p,
		},
	}
}

type serialLegalPersonName struct {
	NameIdentifiers         []*LegalPersonNameId      `json:"nameIdentifier,omitempty"`
	LocalNameIdentifiers    []*LocalLegalPersonNameId `json:"localNameIdentifier,omitempty"`
	PhoneticNameIdentifiers []*LocalLegalPersonNameId `json:"phoneticNameIdentifier,omitempty"`
}

var serialLegalPersonNameFields = map[string]string{
	"nameIdentifier":            "nameIdentifier",
	"nameIdentifiers":           "nameIdentifier",
	"name_identifier":           "nameIdentifier",
	"name_identifiers":          "nameIdentifier",
	"localNameIdentifier":       "localNameIdentifier",
	"localNameIdentifiers":      "localNameIdentifier",
	"local_name_identifier":     "localNameIdentifier",
	"local_name_identifiers":    "localNameIdentifier",
	"phoneticNameIdentifier":    "phoneticNameIdentifier",
	"phoneticNameIdentifiers":   "phoneticNameIdentifier",
	"phonetic_name_identifier":  "phoneticNameIdentifier",
	"phonetic_name_identifiers": "phoneticNameIdentifier",
}

func (l *LegalPersonName) MarshalJSON() ([]byte, error) {
	middle := serialLegalPersonName{
		LocalNameIdentifiers:    l.LocalNameIdentifiers,
		NameIdentifiers:         l.NameIdentifiers,
		PhoneticNameIdentifiers: l.PhoneticNameIdentifiers,
	}
	return json.Marshal(middle)
}

func (l *LegalPersonName) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialLegalPersonNameFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialLegalPersonName{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate legal person values
	l.NameIdentifiers = middle.NameIdentifiers
	l.LocalNameIdentifiers = middle.LocalNameIdentifiers
	l.PhoneticNameIdentifiers = middle.PhoneticNameIdentifiers

	return nil
}

//===========================================================================
// LegalPersonNameID Methods
//===========================================================================

type serialLegalPersonNameID struct {
	LegalPersonName               string                  `json:"legalPersonName,omitempty"`
	LegalPersonNameIdentifierType LegalPersonNameTypeCode `json:"legalPersonNameIdentifierType,omitempty"`
}

var serialLegalPersonNameIDFields = map[string]string{
	"legalPersonName":                   "legalPersonName",
	"legal_person_name":                 "legalPersonName",
	"name":                              "legalPersonName",
	"legalPersonNameIdentifierType":     "legalPersonNameIdentifierType",
	"legal_person_name_identifier_type": "legalPersonNameIdentifierType",
}

func (p *LegalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialLegalPersonNameID{
		LegalPersonName:               p.LegalPersonName,
		LegalPersonNameIdentifierType: p.LegalPersonNameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LegalPersonNameId) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialLegalPersonNameIDFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialLegalPersonNameID{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate legal person values
	p.LegalPersonName = middle.LegalPersonName
	p.LegalPersonNameIdentifierType = middle.LegalPersonNameIdentifierType

	return nil
}

//===========================================================================
// LocalLegalPersonNameID Methods
//===========================================================================

type serialLocalLegalPersonNameID struct {
	LegalPersonName               string                  `json:"legalPersonName,omitempty"`
	LegalPersonNameIdentifierType LegalPersonNameTypeCode `json:"legalPersonNameIdentifierType,omitempty"`
}

var serialLocalLegalPersonNameIDFields = map[string]string{
	"legalPersonName":                   "legalPersonName",
	"legal_person_name":                 "legalPersonName",
	"name":                              "legalPersonName",
	"legalPersonNameIdentifierType":     "legalPersonNameIdentifierType",
	"legal_person_name_identifier_type": "legalPersonNameIdentifierType",
}

func (p *LocalLegalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialLocalLegalPersonNameID{
		LegalPersonName:               p.LegalPersonName,
		LegalPersonNameIdentifierType: p.LegalPersonNameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LocalLegalPersonNameId) UnmarshalJSON(data []byte) (err error) {
	// Perform rekeying operation
	if allowRekeying {
		if data, err = Rekey(data, serialLocalLegalPersonNameIDFields); err != nil {
			return err
		}
	}

	// Unmarshal middle data structure
	middle := serialLocalLegalPersonNameID{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	// Populate legal person values
	p.LegalPersonName = middle.LegalPersonName
	p.LegalPersonNameIdentifierType = middle.LegalPersonNameIdentifierType

	return nil
}
