package models

import (
	"errors"
	"fmt"
	"net"
	"net/mail"

	"github.com/trisacrypto/trisa/pkg/ivms101"
)

// Name searches the IVMS 101 LegalPerson record for the best name to use to represent
// the VASP entity in text. The resolution order is trading name, short name, finally
// falling back on legal name. If there are more than one of each of these types of
// names then the first name is used.
// TODO: also search local names if locale is specified.
func (v *VASP) Name() (string, error) {
	if v.Entity == nil || v.Entity.Name == nil {
		return "", fmt.Errorf("VASP (%s) does not have a valid legal person entity", v.Id)
	}

	names := make([]string, 3)
	for _, name := range v.Entity.Name.NameIdentifiers {
		switch name.LegalPersonNameIdentifierType {
		case ivms101.LegalPersonTrading:
			if names[0] == "" {
				names[0] = name.LegalPersonName
			}
		case ivms101.LegalPersonShort:
			if names[1] == "" {
				names[1] = name.LegalPersonName
			}
		case ivms101.LegalPersonLegal:
			if names[2] == "" {
				names[2] = name.LegalPersonName
			}
		default:
			continue
		}
	}

	for _, name := range names {
		if name != "" {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not find a name for VASP (%s)", v.Id)
}

// Validate checks if the VASP record is complete with all required fields. If partial
// is specified, the validation checks the VASP record as though it hasn't been created.
func (v *VASP) Validate(partial bool) (err error) {
	if !partial && v.Id == "" {
		return errors.New("VASP missing ID field and is not a partial record")
	}

	if v.Id != "" && v.RegisteredDirectory == "" {
		return errors.New("VASP must have a registered directory if it has an ID")
	}

	if v.Entity == nil {
		return errors.New("VASP does not have a legal person entity for KYC operations")
	}

	if err = v.Entity.Validate(); err != nil {
		return err
	}

	if v.Contacts == nil {
		return errors.New("no contact specified on the VASP entity")
	}

	if err = v.Contacts.Validate(); err != nil {
		return err
	}

	if v.CommonName == "" || v.TrisaEndpoint == "" {
		return errors.New("no TRISA endpoint or domain common name")
	}

	host, port, err := net.SplitHostPort(v.TrisaEndpoint)
	if err != nil || host == "" || port == "" {
		return errors.New("could not resolve trisa endpoint host:port")
	}

	if v.VerificationStatus == VerificationState_VERIFIED && v.VerifiedOn == "" {
		return errors.New("VASP is verified but missing verified date")
	}

	if v.VerifiedOn != "" && v.VerificationStatus != VerificationState_VERIFIED {
		return errors.New("VASP has verified on but is not verified")
	}

	if !partial && (v.FirstListed == "" || v.LastUpdated == "") {
		return errors.New("VASP missing first_listed or last_updated timestamps")
	}

	if !partial && len(v.Signature) == 0 {
		return errors.New("VASP missing search signature and is not a partial record")
	}

	return nil
}

// Validate checks is the required contacts are not zero and properly structured.
func (c *Contacts) Validate() (err error) {
	nValid := 0
	if c.Administrative != nil && !c.Administrative.IsZero() {
		if err = c.Administrative.Validate(); err != nil {
			return fmt.Errorf("administrative contact invalid: %s", err)
		}
		nValid++
	}

	if c.Technical != nil && !c.Technical.IsZero() {
		if err = c.Technical.Validate(); err != nil {
			return fmt.Errorf("technical contact invalid: %s", err)
		}
		nValid++
	}

	if c.Legal != nil && !c.Legal.IsZero() {
		if err = c.Legal.Validate(); err != nil {
			return fmt.Errorf("legal contact invalid: %s", err)
		}
		nValid++
	}

	if c.Billing != nil && !c.Billing.IsZero() {
		if err = c.Billing.Validate(); err != nil {
			return fmt.Errorf("billing contact invalid: %s", err)
		}
		nValid++
	}

	if nValid == 0 {
		return errors.New("no contact specified on the VASP entity")
	}
	return nil
}

// Validate checks if a contact record is complete with all required fields.
func (c *Contact) Validate() (err error) {
	// A record must have a name that is longer than 1 character
	if c.Name == "" {
		return errors.New("contact name is required")
	}

	if len(c.Name) < 2 {
		return errors.New("contact name must be longer than one character")
	}

	// If the name is present (e.g. the contact is not zero) then an email that is
	// parseable by RFC 5322 (e.g. by the standard lib mail package).
	if c.Email == "" {
		return errors.New("contact email is required")
	}

	if _, err = mail.ParseAddress(c.Email); err != nil {
		return errors.New("could not parse email address")
	}

	return nil
}

// IsZero returns true if the contact is empty; e.g. it has no name, email, or phone
// number. The Natural Person record on the contact is ignored since this is a deep
// nested structure. The extra field is also ignored since this is side data.
func (c *Contact) IsZero() bool {
	return c.Name == "" && c.Email == "" && c.Phone == ""
}
