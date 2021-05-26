package models

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/trisacrypto/trisa/pkg/ivms101"
)

// Business Category Enumeration Helpers
const (
	BusinessCategoryUnknown       = BusinessCategory_UNKNOWN_ENTITY
	BusinessCategoryPrivate       = BusinessCategory_PRIVATE_ORGANIZATION
	BusinessCategoryGovernment    = BusinessCategory_GOVERNMENT_ENTITY
	BusinessCategoryBusiness      = BusinessCategory_BUSINESS_ENTITY
	BusinessCategoryNonCommercial = BusinessCategory_NON_COMMERCIAL_ENTITY
)

// VASP Category Enumeration Helpers
const (
	VASPCategoryUnknown    = "Unknown"
	VASPCategoryExchange   = "Exchange"
	VASPCategoryDEX        = "DEX"
	VASPCategoryP2P        = "P2P"
	VASPCategoryKiosk      = "Kiosk"
	VASPCategoryCustodian  = "Custodian"
	VASPCategoryOTC        = "OTC"
	VASPCategoryFund       = "Fund"
	VASPCategoryProject    = "Project"
	VASPCategoryGambling   = "Gambling"
	VASPCategoryMiner      = "Miner"
	VASPCategoryMixer      = "Mixer"
	VASPCategoryIndividual = "Individual"
	VASPCategoryOther      = "Other"
)

// ParseBusinessCategory from text representation.
func ParseBusinessCategory(s string) (BusinessCategory, error) {
	s = strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
	code, ok := BusinessCategory_value[s]
	if ok {
		return BusinessCategory(code), nil
	}
	return BusinessCategoryUnknown, fmt.Errorf("could not parse %q into a business category", s)
}

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

	if v.Contacts.Technical == nil &&
		v.Contacts.Billing == nil &&
		v.Contacts.Administrative == nil &&
		v.Contacts.Legal == nil {
		return errors.New("no contact specified on the VASP entity")
	}

	if v.Contacts.Technical != nil && v.Contacts.Technical.Email == "" {
		return errors.New("missing technical contact email")
	}

	if v.Contacts.Billing != nil && v.Contacts.Billing.Email == "" {
		return errors.New("missing billing contact email")
	}

	if v.Contacts.Administrative != nil && v.Contacts.Administrative.Email == "" {
		return errors.New("missing administrative contact email")
	}

	if v.Contacts.Legal != nil && v.Contacts.Legal.Email == "" {
		return errors.New("missing legal contact email")
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
