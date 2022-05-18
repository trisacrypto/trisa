package ivms101

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Short form natural person name type codes.
const (
	NaturalPersonAlias  = NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_ALIA
	NaturalPersonBirth  = NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_BIRT
	NaturalPersonMaiden = NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_MAID
	NaturalPersonLegal  = NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_LEGL
	NaturalPersonMisc   = NaturalPersonNameTypeCode_NATURAL_PERSON_NAME_TYPE_CODE_MISC
)

// Short form legal person name type codes.
const (
	LegalPersonLegal   = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL
	LegalPersonShort   = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_SHRT
	LegalPersonTrading = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_TRAD
)

// Short form address type codes.
const (
	AddressTypeHome       = AddressTypeCode_ADDRESS_TYPE_CODE_HOME
	AddressTypeBusiness   = AddressTypeCode_ADDRESS_TYPE_CODE_BIZZ
	AddressTypeGeographic = AddressTypeCode_ADDRESS_TYPE_CODE_GEOG
)

// Short form national identifier type codes.
const (
	NationalIdentifierARNU = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_ARNU
	NationalIdentifierCCPT = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_CCPT
	NationalIdentifierRAID = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_RAID
	NationalIdentifierDRLC = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_DRLC
	NationalIdentifierFIIN = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_FIIN
	NationalIdentifierTXID = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_TXID
	NationalIdentifierSOCS = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_SOCS
	NationalIdentifierIDCD = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_IDCD
	NationalIdentifierLEIX = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_LEIX
	NationalIdentifierMISC = NationalIdentifierTypeCode_NATIONAL_IDENTIFIER_TYPE_CODE_MISC
)

//
// Natural Person Enums and Codes
//

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (n NaturalPersonNameTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(n.String(), "NATURAL_PERSON_NAME_TYPE_CODE_")
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (n *NaturalPersonNameTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse NaturalPersonNameTypeCode from value")
	}
	s = "NATURAL_PERSON_NAME_TYPE_CODE_" + strings.ToUpper(s)
	code, ok := NaturalPersonNameTypeCode_value[s]
	if !ok {
		return errors.New("invalid NaturalPersonNameTypeCode alias")
	}
	*n = NaturalPersonNameTypeCode(code)
	return nil
}

func (p *NaturalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := struct {
		PrimaryIdentifier   string                    `json:"primaryIdentifier,omitempty"`
		SecondaryIdentifier string                    `json:"secondaryIdentifier,omitempty"`
		NameIdentifierType  NaturalPersonNameTypeCode `json:"nameIdentifierType,omitempty"`
	}{
		PrimaryIdentifier:   p.PrimaryIdentifier,
		SecondaryIdentifier: p.SecondaryIdentifier,
		NameIdentifierType:  p.NameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *NaturalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := make(map[string]string)
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}
	var code NaturalPersonNameTypeCode
	t := fmt.Sprintf("\"%s\"", middle["nameIdentifierType"])
	err := code.UnmarshalJSON([]byte(t))
	if err != nil {
		return err
	}
	i := NaturalPersonNameId{
		PrimaryIdentifier:   middle["primaryIdentifier"],
		SecondaryIdentifier: middle["secondaryIdentifier"],
		NameIdentifierType:  code,
	}
	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

func (p *LocalNaturalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := struct {
		PrimaryIdentifier   string                    `json:"primaryIdentifier,omitempty"`
		SecondaryIdentifier string                    `json:"secondaryIdentifier,omitempty"`
		NameIdentifierType  NaturalPersonNameTypeCode `json:"nameIdentifierType,omitempty"`
	}{
		PrimaryIdentifier:   p.PrimaryIdentifier,
		SecondaryIdentifier: p.SecondaryIdentifier,
		NameIdentifierType:  p.NameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LocalNaturalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := make(map[string]string)
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}
	var code NaturalPersonNameTypeCode
	t := fmt.Sprintf("\"%s\"", middle["nameIdentifierType"])
	err := code.UnmarshalJSON([]byte(t))
	if err != nil {
		return err
	}
	i := LocalNaturalPersonNameId{
		PrimaryIdentifier:   middle["primaryIdentifier"],
		SecondaryIdentifier: middle["secondaryIdentifier"],
		NameIdentifierType:  code,
	}
	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

func (n *NaturalPersonName) MarshalJSON() ([]byte, error) {
	middle := struct {
		LocalNameIdentifiers    []*LocalNaturalPersonNameId `json:"localNameIdentifier,omitempty"`
		NameIdentifiers         []*NaturalPersonNameId      `json:"nameIdentifier,omitempty"`
		PhoneticNameIdentifiers []*LocalNaturalPersonNameId `json:"phoneticNameIdentifier,omitempty"`
	}{
		LocalNameIdentifiers:    n.LocalNameIdentifiers,
		NameIdentifiers:         n.NameIdentifiers,
		PhoneticNameIdentifiers: n.PhoneticNameIdentifiers,
	}
	return json.Marshal(middle)
}

func (n *NaturalPersonName) UnmarshalJSON(data []byte) error {
	middle := struct {
		LocalNameIdentifiers    []*LocalNaturalPersonNameId `json:"localNameIdentifier"`
		NameIdentifiers         []*NaturalPersonNameId      `json:"nameIdentifier"`
		PhoneticNameIdentifiers []*LocalNaturalPersonNameId `json:"phoneticNameIdentifier"`
	}{}
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

// TODO
// func (a *NaturalPerson) MarshalJSON() ([]byte, error) {
// 	return nil, nil
// }

// TODO
// func (a *NaturalPerson) UnmarshalJSON(data []byte) error {
// 	return nil
// }

//
// Legal Person Enums and Codes
//

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (l LegalPersonNameTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(l.String(), "LEGAL_PERSON_NAME_TYPE_CODE_")
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (l *LegalPersonNameTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse LegalPersonNameTypeCode from value")
	}
	s = "LEGAL_PERSON_NAME_TYPE_CODE_" + strings.ToUpper(s)
	code, ok := LegalPersonNameTypeCode_value[s]
	if !ok {
		return errors.New("invalid LegalPersonNameTypeCode alias")
	}
	*l = LegalPersonNameTypeCode(code)
	return nil
}

func (p *LegalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := struct {
		LegalPersonName               string                  `json:"legalPersonName,omitempty"`
		LegalPersonNameIdentifierType LegalPersonNameTypeCode `json:"legalPersonNameIdentifierType,omitempty"`
	}{
		LegalPersonName:               p.LegalPersonName,
		LegalPersonNameIdentifierType: p.LegalPersonNameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LegalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := make(map[string]string)
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}
	var code LegalPersonNameTypeCode
	t := fmt.Sprintf("\"%s\"", middle["legalPersonNameIdentifierType"])
	err := code.UnmarshalJSON([]byte(t))
	if err != nil {
		return err
	}
	i := LegalPersonNameId{
		LegalPersonName:               middle["legalPersonName"],
		LegalPersonNameIdentifierType: code,
	}
	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

func (p *LocalLegalPersonNameId) MarshalJSON() ([]byte, error) {
	middle := struct {
		LegalPersonName               string                  `json:"legalPersonName,omitempty"`
		LegalPersonNameIdentifierType LegalPersonNameTypeCode `json:"legalPersonNameIdentifierType,omitempty"`
	}{
		LegalPersonName:               p.LegalPersonName,
		LegalPersonNameIdentifierType: p.LegalPersonNameIdentifierType,
	}
	return json.Marshal(middle)
}

func (p *LocalLegalPersonNameId) UnmarshalJSON(data []byte) error {
	middle := make(map[string]string)
	if err := json.Unmarshal(data, &middle); err != nil {
		return err
	}
	var code LegalPersonNameTypeCode
	t := fmt.Sprintf("\"%s\"", middle["legalPersonNameIdentifierType"])
	err := code.UnmarshalJSON([]byte(t))
	if err != nil {
		return err
	}
	i := LocalLegalPersonNameId{
		LegalPersonName:               middle["legalPersonName"],
		LegalPersonNameIdentifierType: code,
	}
	// TODO warning: assignment copies lock value to *p
	*p = i
	return nil
}

func (l *LegalPersonName) MarshalJSON() ([]byte, error) {
	middle := struct {
		LocalNameIdentifiers    []*LocalLegalPersonNameId `json:"localNameIdentifier,omitempty"`
		NameIdentifiers         []*LegalPersonNameId      `json:"nameIdentifier,omitempty"`
		PhoneticNameIdentifiers []*LocalLegalPersonNameId `json:"phoneticNameIdentifier,omitempty"`
	}{
		LocalNameIdentifiers:    l.LocalNameIdentifiers,
		NameIdentifiers:         l.NameIdentifiers,
		PhoneticNameIdentifiers: l.PhoneticNameIdentifiers,
	}
	return json.Marshal(middle)
}

func (l *LegalPersonName) UnmarshalJSON(data []byte) error {
	middle := struct {
		LocalNameIdentifiers    []*LocalLegalPersonNameId `json:"localNameIdentifier"`
		NameIdentifiers         []*LegalPersonNameId      `json:"nameIdentifier"`
		PhoneticNameIdentifiers []*LocalLegalPersonNameId `json:"phoneticNameIdentifier"`
	}{}
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

// TODO
// func (e *LegalPerson) MarshalJSON() ([]byte, error) {
// 	return nil, nil
// }

// TODO
// func (e *LegalPerson) UnmarshalJSON(data []byte) error {
// 	return nil
// }

//
// National Identifier Enums and Codes
//

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (i NationalIdentifierTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(i.String(), "NATIONAL_IDENTIFIER_TYPE_CODE_")
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (i *NationalIdentifierTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse NationalIdentifierTypeCode from value")
	}
	s = "NATIONAL_IDENTIFIER_TYPE_CODE_" + strings.ToUpper(s)
	code, ok := NationalIdentifierTypeCode_value[s]
	if !ok {
		return errors.New("invalid NationalIdentifierTypeCode alias")
	}
	*i = NationalIdentifierTypeCode(code)
	return nil
}

//
// Geographic Address Enums and Codes
//

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (a AddressTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(a.String(), "ADDRESS_TYPE_CODE_")
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (a *AddressTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse AddressTypeCode from value")
	}
	s = "ADDRESS_TYPE_CODE_" + strings.ToUpper(s)
	code, ok := AddressTypeCode_value[s]
	if !ok {
		return errors.New("invalid AddressTypeCode alias")
	}
	*a = AddressTypeCode(code)
	return nil
}

func (a *Address) MarshalJSON() ([]byte, error) {
	middle := struct {
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
	}{
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
	middle := struct {
		AddressType        AddressTypeCode `json:"addressType"`
		Department         string          `json:"department"`
		SubDepartment      string          `json:"subDepartment"`
		StreetName         string          `json:"streetName"`
		BuildingNumber     string          `json:"buildingNumber"`
		BuildingName       string          `json:"buildingName"`
		Floor              string          `json:"floor"`
		PostBox            string          `json:"postBox"`
		Room               string          `json:"room"`
		PostCode           string          `json:"postCode"`
		TownName           string          `json:"townName"`
		TownLocationName   string          `json:"townLocationName"`
		DistrictName       string          `json:"districtName"`
		CountrySubDivision string          `json:"countrySubDivision"`
		AddressLine        []string        `json:"addressLine"`
		Country            string          `json:"country"`
	}{}
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
