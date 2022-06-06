package ivms101

import (
	"encoding/json"
	"errors"
)

//
// Person JSON
//

type serialPerson struct {
	NaturalPerson *NaturalPerson `json:"naturalPerson,omitempty"`
	LegalPerson   *LegalPerson   `json:"legalPerson,omitempty"`
}

func (p *Person) MarshalJSON() ([]byte, error) {
	middle := serialPerson{
		NaturalPerson: p.GetNaturalPerson(),
		LegalPerson:   p.GetLegalPerson(),
	}
	return json.Marshal(middle)
}

func (p *Person) UnmarshalJSON(data []byte) error {
	middle := serialPerson{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	if middle.NaturalPerson != nil && middle.LegalPerson != nil {
		return errors.New("person object cannot be both a natural and legal person")
	}

	i := Person{}
	switch {
	case middle.NaturalPerson != nil:
		i.Person = &Person_NaturalPerson{
			NaturalPerson: middle.NaturalPerson,
		}
	case middle.LegalPerson != nil:
		i.Person = &Person_LegalPerson{
			LegalPerson: middle.LegalPerson,
		}
	}

	*p = i
	return nil
}

//
// NaturalPerson JSON
//

type serialNaturalPerson struct {
	Name               *NaturalPersonName      `json:"name,omitempty"`
	Address            []*Address              `json:"geographicAddress,omitempty"`
	Identification     *NationalIdentification `json:"nationalIdentification,omitempty"`
	CustomerID         string                  `json:"customerIdentification,omitempty"`
	DOB                *DateAndPlaceOfBirth    `json:"dateAndPlaceOfBirth,omitempty"`
	CountryOfResidence string                  `json:"countryOfResidence,omitempty"`
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

func (n *NaturalPerson) UnmarshalJSON(data []byte) error {
	middle := serialNaturalPerson{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := NaturalPerson{
		Name:                   middle.Name,
		GeographicAddresses:    middle.Address,
		NationalIdentification: middle.Identification,
		CustomerIdentification: middle.CustomerID,
		DateAndPlaceOfBirth:    middle.DOB,
		CountryOfResidence:     middle.CountryOfResidence,
	}

	// TODO warning: assignment copies lock value to *n
	*n = i
	return nil
}

//
// NaturalPersonName and NaturalPersonNameIdentifiers JSON
//

type serialNaturalPersonName struct {
	LocalNameIdentifiers    []*LocalNaturalPersonNameId `json:"localNameIdentifier,omitempty"`
	NameIdentifiers         []*NaturalPersonNameId      `json:"nameIdentifier,omitempty"`
	PhoneticNameIdentifiers []*LocalNaturalPersonNameId `json:"phoneticNameIdentifier,omitempty"`
}

func (n *NaturalPersonName) MarshalJSON() ([]byte, error) {
	middle := serialNaturalPersonName{
		LocalNameIdentifiers:    n.LocalNameIdentifiers,
		NameIdentifiers:         n.NameIdentifiers,
		PhoneticNameIdentifiers: n.PhoneticNameIdentifiers,
	}
	return json.Marshal(middle)
}

func (n *NaturalPersonName) UnmarshalJSON(data []byte) error {
	middle := serialNaturalPersonName{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := NaturalPersonName{
		NameIdentifiers:         middle.NameIdentifiers,
		LocalNameIdentifiers:    middle.LocalNameIdentifiers,
		PhoneticNameIdentifiers: middle.PhoneticNameIdentifiers,
	}

	// TODO warning: assignment copies lock value to *n
	*n = i
	return nil
}

type serialNaturalPersonNameId struct {
	PrimaryIdentifier   string                    `json:"primaryIdentifier,omitempty"`
	SecondaryIdentifier string                    `json:"secondaryIdentifier,omitempty"`
	NameIdentifierType  NaturalPersonNameTypeCode `json:"nameIdentifierType,omitempty"`
}

func (p *NaturalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialNaturalPersonNameId{
		PrimaryIdentifier:   p.PrimaryIdentifier,
		SecondaryIdentifier: p.SecondaryIdentifier,
		NameIdentifierType:  p.NameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *NaturalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := serialNaturalPersonNameId{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := NaturalPersonNameId{
		PrimaryIdentifier:   middle.PrimaryIdentifier,
		SecondaryIdentifier: middle.SecondaryIdentifier,
		NameIdentifierType:  middle.NameIdentifierType,
	}

	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

type serialLocalNaturalPersonNameId struct {
	PrimaryIdentifier   string                    `json:"primaryIdentifier,omitempty"`
	SecondaryIdentifier string                    `json:"secondaryIdentifier,omitempty"`
	NameIdentifierType  NaturalPersonNameTypeCode `json:"nameIdentifierType,omitempty"`
}

func (p *LocalNaturalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialLocalNaturalPersonNameId{
		PrimaryIdentifier:   p.PrimaryIdentifier,
		SecondaryIdentifier: p.SecondaryIdentifier,
		NameIdentifierType:  p.NameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LocalNaturalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := serialLocalNaturalPersonNameId{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := LocalNaturalPersonNameId{
		PrimaryIdentifier:   middle.PrimaryIdentifier,
		SecondaryIdentifier: middle.SecondaryIdentifier,
		NameIdentifierType:  middle.NameIdentifierType,
	}

	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

//
// Address JSON
//

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

func (a *Address) UnmarshalJSON(data []byte) error {
	middle := serialAddress{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := Address{
		AddressType:        middle.AddressType,
		Department:         middle.Department,
		SubDepartment:      middle.SubDepartment,
		StreetName:         middle.StreetName,
		BuildingNumber:     middle.BuildingNumber,
		BuildingName:       middle.BuildingName,
		Floor:              middle.Floor,
		PostBox:            middle.PostBox,
		Room:               middle.Room,
		PostCode:           middle.PostCode,
		TownName:           middle.TownName,
		TownLocationName:   middle.TownLocationName,
		DistrictName:       middle.DistrictName,
		CountrySubDivision: middle.CountrySubDivision,
		AddressLine:        middle.AddressLine,
		Country:            middle.Country,
	}

	// TODO warning: assignment copies lock value to *a
	*a = i
	return nil
}

//
// DateAndPlaceOfBirth JSON
//

type serialDateAndPlaceOfBirth struct {
	DateOfBirth  string `json:"dateOfBirth,omitempty"`
	PlaceOfBirth string `json:"placeOfBirth,omitempty"`
}

func (d *DateAndPlaceOfBirth) MarshalJSON() ([]byte, error) {
	middle := serialDateAndPlaceOfBirth{
		DateOfBirth:  d.DateOfBirth,
		PlaceOfBirth: d.PlaceOfBirth,
	}
	return json.Marshal(middle)
}

func (d *DateAndPlaceOfBirth) UnmarshalJSON(data []byte) error {
	middle := serialDateAndPlaceOfBirth{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := DateAndPlaceOfBirth{
		DateOfBirth:  middle.DateOfBirth,
		PlaceOfBirth: middle.PlaceOfBirth,
	}

	// TODO warning: assignment copies lock value to *d
	*d = i
	return nil
}

//
// NationalIdentification JSON
//

type serialNationalIdentification struct {
	NationalIdentifier     string                     `json:"nationalIdentifier,omitempty"`
	NationalIdentifierType NationalIdentifierTypeCode `json:"nationalIdentifierType,omitempty"`
	CountryOfIssue         string                     `json:"countryOfIssue,omitempty"`
	RegistrationAuthority  string                     `json:"registrationAuthority,omitempty"`
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

func (n *NationalIdentification) UnmarshalJSON(data []byte) error {
	middle := serialNationalIdentification{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := NationalIdentification{
		NationalIdentifier:     middle.NationalIdentifier,
		NationalIdentifierType: middle.NationalIdentifierType,
		CountryOfIssue:         middle.CountryOfIssue,
		RegistrationAuthority:  middle.RegistrationAuthority,
	}

	// TODO warning: assignment copies lock value to *n
	*n = i
	return nil
}

//
// LegalPerson JSON
//

type serialLegalPerson struct {
	Name                   *LegalPersonName        `json:"name,omitempty"`
	Address                []*Address              `json:"geographicAddress,omitempty"`
	CustomerNumber         string                  `json:"customerNumber,omitempty"`
	NationalIdentification *NationalIdentification `json:"nationalIdentification,omitempty"`
	CountryOfRegistration  string                  `json:"countryOfRegistration,omitempty"`
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
	middle := serialLegalPerson{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := LegalPerson{
		Name:                   middle.Name,
		GeographicAddresses:    middle.Address,
		CustomerNumber:         middle.CustomerNumber,
		NationalIdentification: middle.NationalIdentification,
		CountryOfRegistration:  middle.CountryOfRegistration,
	}

	// TODO warning: assignment copies lock value to *n
	*l = i
	return nil
}

//
// LegalPersonName and LegalPersonNameIdentifiers JSON
//

type serialLegalPersonName struct {
	LocalNameIdentifiers    []*LocalLegalPersonNameId `json:"localNameIdentifier,omitempty"`
	NameIdentifiers         []*LegalPersonNameId      `json:"nameIdentifier,omitempty"`
	PhoneticNameIdentifiers []*LocalLegalPersonNameId `json:"phoneticNameIdentifier,omitempty"`
}

func (l *LegalPersonName) MarshalJSON() ([]byte, error) {
	middle := serialLegalPersonName{
		LocalNameIdentifiers:    l.LocalNameIdentifiers,
		NameIdentifiers:         l.NameIdentifiers,
		PhoneticNameIdentifiers: l.PhoneticNameIdentifiers,
	}
	return json.Marshal(middle)
}

func (l *LegalPersonName) UnmarshalJSON(data []byte) error {
	middle := serialLegalPersonName{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := LegalPersonName{
		NameIdentifiers:         middle.NameIdentifiers,
		LocalNameIdentifiers:    middle.LocalNameIdentifiers,
		PhoneticNameIdentifiers: middle.PhoneticNameIdentifiers,
	}

	// TODO warning: assignment copies lock value to *l
	*l = i
	return nil
}

type serialLegalPersonNameId struct {
	LegalPersonName               string                  `json:"legalPersonName,omitempty"`
	LegalPersonNameIdentifierType LegalPersonNameTypeCode `json:"legalPersonNameIdentifierType,omitempty"`
}

func (p *LegalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialLegalPersonNameId{
		LegalPersonName:               p.LegalPersonName,
		LegalPersonNameIdentifierType: p.LegalPersonNameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LegalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := serialLegalPersonNameId{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := LegalPersonNameId{
		LegalPersonName:               middle.LegalPersonName,
		LegalPersonNameIdentifierType: middle.LegalPersonNameIdentifierType,
	}

	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

type serialLocalLegalPersonNameId struct {
	LegalPersonName               string                  `json:"legalPersonName,omitempty"`
	LegalPersonNameIdentifierType LegalPersonNameTypeCode `json:"legalPersonNameIdentifierType,omitempty"`
}

func (p *LocalLegalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := serialLocalLegalPersonNameId{
		LegalPersonName:               p.LegalPersonName,
		LegalPersonNameIdentifierType: p.LegalPersonNameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LocalLegalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := serialLocalLegalPersonNameId{}
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}

	i := LocalLegalPersonNameId{
		LegalPersonName:               middle.LegalPersonName,
		LegalPersonNameIdentifierType: middle.LegalPersonNameIdentifierType,
	}

	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}
