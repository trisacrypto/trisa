package ivms101

import "errors"

// Standard error values for error type checking
var (
	ErrNoNaturalPersonNameIdentifiers           = errors.New("one or more natural person name identifiers is required")
	ErrInvalidNaturalPersonName                 = errors.New("natural person name required with max length 100 chars")
	ErrInvalidNaturalPersonNameTypeCode         = errors.New("invalid natural person name type code")
	ErrParseNaturalPersonNameTypeCode           = errors.New("could not parse natural person name type code from value")
	ErrNoLegalPersonNameIdentifiers             = errors.New("one or more legal person name identifiers is required")
	ErrInvalidLegalPersonName                   = errors.New("legal person name required with max length 100 chars")
	ErrInvalidLegalPersonNameTypeCode           = errors.New("invalid legal person name type code")
	ErrParseLegalPersonNameTypeCode             = errors.New("could not parse legal person name type code from value")
	ErrLegalNamesPresent                        = errors.New("at least one name identifier must have a LEGL name identifier type")
	ErrInvalidCustomerNumber                    = errors.New("customer number can be at most 50 chars")
	ErrInvalidCustomerIdentification            = errors.New("customer identification can be at most 50 chars")
	ErrInvalidCountryCode                       = errors.New("invalid ISO-3166-1 alpha-2 country code")
	ErrValidNationalIdentifierLegalPerson       = errors.New("a legal person must have a national identifier of type RAID, MISC, LEIX, or TXID")
	ErrInvalidLEI                               = errors.New("national identifier required with max length 35")
	ErrInvalidNationalIdentifierTypeCode        = errors.New("invalid national identifier type code")
	ErrParseNationalIdentifierTypeCode          = errors.New("could not parse national identifier type code from value")
	ErrCompleteNationalIdentifierCountry        = errors.New("a legal person must not have a value for country if identifier type is not LEIX")
	ErrCompleteNationalIdentifierAuthorityEmpty = errors.New("a legal person must have a value for registration authority if identifier type is not LEIX")
	ErrCompleteNationalIdentifierAuthority      = errors.New("a legal person must not have a value for registration authority if identifier type is LEIX")
	ErrInvalidDateOfBirth                       = errors.New("date of birth must be a valid date in YYYY-MM-DD format")
	ErrInvalidPlaceOfBirth                      = errors.New("place of birth required with at most 70 characters")
	ErrDateInPast                               = errors.New("date of birth must be a historic date, prior to current date")
	ErrValidAddress                             = errors.New("address must have at least one address line or street name + building name or number")
	ErrInvalidAddressTypeCode                   = errors.New("invalid address type code")
	ErrParseAddressTypeCode                     = errors.New("could not parse address type code from value")
	ErrInvalidAddressLines                      = errors.New("an address can contain at most 7 address lines")
	ErrInvalidTransliterationMethodCode         = errors.New("invalid transliteration method code")
	ErrParseTransliterationMethodCode           = errors.New("could not parse transliteration method code from value")
)

// Wraps one error with another error for better error type checking.
type WrappedError struct {
	error error
	cause error
}

func (e *WrappedError) Error() string { return e.error.Error() }

func (e *WrappedError) Cause() error { return e.cause }

func (e *WrappedError) Unwrap() []error { return []error{e.error, e.cause} }

func Wrap(err, cause error) error {
	return &WrappedError{error: err, cause: cause}
}
