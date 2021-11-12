/*
Package ivms101 extends the Go protocol buffers generated by the ivms101 protobuf
package with JSON loading utilities, validation helpers, short constants, etc.
*/
package ivms101

import (
	"strings"
	"time"
)

// Person converts a NaturalPerson into a Person protobuf message type.
func (p *NaturalPerson) Person() *Person {
	return &Person{
		Person: &Person_NaturalPerson{
			NaturalPerson: p,
		},
	}
}

// Person converts a LegalPerson into a Person protobuf message type.
func (p *LegalPerson) Person() *Person {
	return &Person{
		Person: &Person_LegalPerson{
			LegalPerson: p,
		},
	}
}

// Validate the IVMS101 constraints for a natural person data definition. ON the first
// invalid constraint found an error is returned.  No error is returned for valid data.
func (p *NaturalPerson) Validate() (err error) {
	// Constraint: required ValidNaturalPersonName
	if p.Name == nil {
		return ErrNoNaturalPersonNameIdentifiers
	}
	if err = p.Name.Validate(); err != nil {
		return err
	}

	// Constraint: ValidAddresses
	for _, addr := range p.GeographicAddresses {
		if err = addr.Validate(); err != nil {
			return err
		}
	}

	// Constraint: Optional ValidNationalIdentification
	if p.NationalIdentification != nil {
		if err = p.NationalIdentification.Validate(); err != nil {
			return err
		}
	}

	// Constraint: optional Max50Text
	if len(p.CustomerIdentification) > 50 {
		return ErrInvalidCustomerIdentification
	}

	// Constraint: Optional Valid DateAndPlaceOfBirth
	if p.DateAndPlaceOfBirth != nil {
		if err = p.DateAndPlaceOfBirth.Validate(); err != nil {
			return err
		}
	}

	// Constraint: Optional ISO-3166-1 alpha-2 codes or XX
	if p.CountryOfResidence != "" {
		// TODO: ensure the code is valid; for now just checking length
		if len(p.CountryOfResidence) != 2 {
			return ErrInvalidCountryCode
		}

		// Ensure that country code is all upper case
		p.CountryOfResidence = strings.ToUpper(p.CountryOfResidence)
	}

	return nil
}

// Validate the IVMS101 constraints for a legal person data definition. On the first
// invalid constraint found an error is returned. No error is returned for valid data.
func (p *LegalPerson) Validate() (err error) {
	// Constraint: ValidLegalPersonName
	// Constraint: LegalNamePresentLegalPerson
	if p.Name == nil {
		return ErrNoLegalPersonNameIdentifiers
	}
	if err = p.Name.Validate(); err != nil {
		return err
	}

	// Constraint: ValidAddresses
	for _, addr := range p.GeographicAddresses {
		if err = addr.Validate(); err != nil {
			return err
		}
	}

	// Constraint: Optional Max50Text Datatype
	if p.CustomerNumber != "" && len(p.CustomerNumber) > 50 {
		return ErrInvalidCustomerNumber
	}

	// Constraint: Optional ValidNationalIdentification
	if p.NationalIdentification != nil {
		if err = p.NationalIdentification.Validate(); err != nil {
			return err
		}

		// Constraint: ValidNationalIdentifierLegalPerson
		if !(p.NationalIdentification.NationalIdentifierType == NationalIdentifierRAID ||
			p.NationalIdentification.NationalIdentifierType == NationalIdentifierMISC ||
			p.NationalIdentification.NationalIdentifierType == NationalIdentifierLEIX ||
			p.NationalIdentification.NationalIdentifierType == NationalIdentifierTXID) {
			return ErrValidNationalIdentifierLegalPerson
		}

		// Constraint: CompleteNationalIdentifierLegalPerson
		if p.NationalIdentification.NationalIdentifierType != NationalIdentifierLEIX {
			if p.NationalIdentification.CountryOfIssue != "" || p.NationalIdentification.RegistrationAuthority == "" {
				return ErrCompleteNationalIdentifierLegalPerson
			}
		}

	}

	// Constraint: Optional ISO-3166-1 alpha-2 codes or XX
	if p.CountryOfRegistration != "" {
		// TODO: ensure the code is valid; for now just checking length
		if len(p.CountryOfRegistration) != 2 {
			return ErrInvalidCountryCode
		}

		// Ensure that country code is all upper case
		p.CountryOfRegistration = strings.ToUpper(p.CountryOfRegistration)
	}

	return nil
}

// Validate the IVMS101 constraints for natural person name
func (n *NaturalPersonName) Validate() (err error) {
	// Constraint one or more
	if len(n.NameIdentifiers) < 1 {
		return ErrNoNaturalPersonNameIdentifiers
	}

	// Constraint: valid name identifiers
	var legalNames int
	for _, name := range n.NameIdentifiers {
		if err = name.Validate(); err != nil {
			return err
		}

		if name.NameIdentifierType == NaturalPersonLegal {
			legalNames++
		}
	}

	// Constraint: LegalNamePresent
	if legalNames == 0 {
		return ErrLegalNamesPresent
	}

	// Constraint: valid local name identifiers
	for _, name := range n.LocalNameIdentifiers {
		if err = name.Validate(); err != nil {
			return err
		}
	}

	// Constraint: valid phonetic name identifiers
	for _, name := range n.PhoneticNameIdentifiers {
		if err = name.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate the IVMS101 constraints for natural person name identifiers
func (n *NaturalPersonNameId) Validate() (err error) {
	if n.PrimaryIdentifier == "" || len(n.PrimaryIdentifier) > 100 {
		return ErrInvalidNaturalPersonName
	}

	if len(n.SecondaryIdentifier) > 100 {
		return ErrInvalidNaturalPersonName
	}

	typeCode := int32(n.NameIdentifierType)
	if _, ok := NaturalPersonNameTypeCode_name[typeCode]; !ok {
		return ErrInvalidNaturalPersonNameTypeCode
	}

	return nil
}

// Validate the IVMS101 constraints for local natural person anme identifiers
func (n *LocalNaturalPersonNameId) Validate() (err error) {
	if n.PrimaryIdentifier == "" || len(n.PrimaryIdentifier) > 100 {
		return ErrInvalidNaturalPersonName
	}

	if len(n.SecondaryIdentifier) > 100 {
		return ErrInvalidNaturalPersonName
	}

	typeCode := int32(n.NameIdentifierType)
	if _, ok := NaturalPersonNameTypeCode_name[typeCode]; !ok {
		return ErrInvalidNaturalPersonNameTypeCode
	}

	return nil
}

// Validate the IVMS101 constraints for legal person name.
func (n *LegalPersonName) Validate() (err error) {
	// Constraint: one or more
	if len(n.NameIdentifiers) < 1 {
		return ErrNoLegalPersonNameIdentifiers
	}

	// Constraint: valid name identifiers
	var legalNames int
	for _, name := range n.NameIdentifiers {
		if err = name.Validate(); err != nil {
			return err
		}

		if name.LegalPersonNameIdentifierType == LegalPersonLegal {
			legalNames++
		}
	}

	// Constraint: LegalNamePresent
	if legalNames == 0 {
		return ErrLegalNamesPresent
	}

	// Constraint: valid local name identifiers
	for _, name := range n.LocalNameIdentifiers {
		if err = name.Validate(); err != nil {
			return err
		}
	}

	// Constraint: valid phonetic name identifiers
	for _, name := range n.PhoneticNameIdentifiers {
		if err = name.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate the IVMS101 constraints for legal person name identifier
func (n *LegalPersonNameId) Validate() (err error) {
	if n.LegalPersonName == "" || len(n.LegalPersonName) > 100 {
		return ErrInvalidLegalPersonName
	}

	typeCode := int32(n.LegalPersonNameIdentifierType)
	if _, ok := LegalPersonNameTypeCode_name[typeCode]; !ok {
		return ErrInvalidLegalPersonNameTypeCode
	}

	return nil
}

// Validate the IVMS101 constraints for local legal person name identifier
func (n *LocalLegalPersonNameId) Validate() (err error) {
	if n.LegalPersonName == "" || len(n.LegalPersonName) > 100 {
		return ErrInvalidLegalPersonName
	}

	typeCode := int32(n.LegalPersonNameIdentifierType)
	if _, ok := LegalPersonNameTypeCode_name[typeCode]; !ok {
		return ErrInvalidLegalPersonNameTypeCode
	}

	return nil
}

// Validate the IVMS101 constraints for a geographic address
func (a *Address) Validate() (err error) {
	// Constraint: valid required address type code
	typeCode := int32(a.AddressType)
	if _, ok := AddressTypeCode_name[typeCode]; !ok {
		return ErrInvalidAddressTypeCode
	}

	// TODO: validate optional max length constraints
	// Constraint: at most 7 address lines
	if len(a.AddressLine) > 7 {
		return ErrInvalidAddressLines
	}

	// Constraint: ValidAddress
	if len(a.AddressLine) == 0 && (a.StreetName == "" && (a.BuildingName == "" || a.BuildingNumber == "")) {
		return ErrValidAddress
	}

	// Constraint: required valid country code
	// TODO: validate ISO-3166-1 alpha-2 country code
	if a.Country == "" || len(a.Country) != 2 {
		return ErrInvalidCountryCode
	}
	a.Country = strings.ToUpper(a.Country)

	return nil
}

// Validate the IVMS101 constraints for a national identification
func (id *NationalIdentification) Validate() (err error) {
	// TODO: Constraint ValidLEI
	// Constraint: required Max35Text datatype
	if id.NationalIdentifier == "" || len(id.NationalIdentifier) > 35 {
		return ErrInvalidLEI
	}

	// Constraint: required valid national identifier type code
	typeCode := int32(id.NationalIdentifierType)
	if _, ok := NationalIdentifierTypeCode_name[typeCode]; !ok {
		return ErrInvalidNationalIdentifierTypeCode
	}

	// Constraint: valid country code
	if id.CountryOfIssue != "" {
		// TODO: validate ISO-3166-1 alpha-2 country code
		if len(id.CountryOfIssue) != 2 {
			return ErrInvalidCountryCode
		}
		id.CountryOfIssue = strings.ToUpper(id.CountryOfIssue)
	}

	// TODO: Contraint authority in GLEIF Registration authorities list
	return nil
}

// Validate the IVMS101 constraints for date and place of birth
func (d *DateAndPlaceOfBirth) Validate() (err error) {
	// Constraint: require valid date
	if d.DateOfBirth == "" {
		return ErrInvalidDateOfBirth
	}

	var date time.Time
	if date, err = time.Parse("2006-01-02", d.DateOfBirth); err != nil {
		return ErrInvalidDateOfBirth
	}

	if d.PlaceOfBirth == "" || len(d.PlaceOfBirth) > 70 {
		return ErrInvalidPlaceOfBirth
	}

	// Constraint: DateInPast
	if date.After(time.Now()) {
		return ErrDateInPast
	}

	return nil
}
