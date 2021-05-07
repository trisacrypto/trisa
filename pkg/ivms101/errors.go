package ivms101

import "errors"

// Standard error values for error type checking
var (
	ErrNoNaturalPersonNameIdentifiers        = errors.New("one or more natural person name identifiers is required")
	ErrInvalidNaturalPersonName              = errors.New("natural person name required with max length 100 chars")
	ErrInvalidNaturalPersonNameTypeCode      = errors.New("invalid natural person name type code")
	ErrNoLegalPersonNameIdentifiers          = errors.New("one or more legal person name identifiers is required")
	ErrInvalidLegalPersonName                = errors.New("legal person name required with max length 100 chars")
	ErrInvalidLegalPersonNameTypeCode        = errors.New("invalid legal person name type code")
	ErrLegalNamesPresent                     = errors.New("at least one name identifier must have a LEGL name identifier type")
	ErrInvalidCustomerNumber                 = errors.New("customer number can be at most 50 chars")
	ErrInvalidCustomerIdentification         = errors.New("customer identification can be at most 50 chars")
	ErrInvalidCountryCode                    = errors.New("invalid ISO-3166-1 alpha-2 country code")
	ErrValidNationalIdentifierLegalPerson    = errors.New("a legal person must have a national identifier of type RAID, MISC, LEIX, or TXID")
	ErrInvalidLEI                            = errors.New("national identifier required with max length 35")
	ErrInvalidNationalIdentifierTypeCode     = errors.New("invalid national identifier type code")
	ErrCompleteNationalIdentifierLegalPerson = errors.New("a legal person must not have a value for country and must have value for registration authority if identifier type is not LEIX")
	ErrInvalidDateOfBirth                    = errors.New("date of birth must be a valid date in YYYY-MM-DD format")
	ErrInvalidPlaceOfBirth                   = errors.New("place of birth required with at most 70 characters")
	ErrDateInPast                            = errors.New("date of birth must be a historic date, prior to current date")
	ErrValidAddress                          = errors.New("address must have at least one address line or street name + building name or number")
	ErrInvalidAddressTypeCode                = errors.New("invalid address type code")
	ErrInvalidAddressLines                   = errors.New("an address can contain at most 7 address lines")
)
