package ivms101

import (
	"errors"
	"fmt"
	"strings"
)

// Standard error values for error type checking
var (
	ErrNoLegalPersonNameIdentifiers             = errors.New("one or more legal person name identifiers is required")
	ErrInvalidLegalPersonName                   = errors.New("legal person name required with max length 100 chars")
	ErrInvalidCustomerNumber                    = errors.New("customer number can be at most 50 chars")
	ErrInvalidCountryCode                       = errors.New("invalid ISO-3166-1 alpha-2 country code")
	ErrValidNationalIdentifierLegalPerson       = errors.New("a legal person must have a national identifier of type RAID, MISC, LEIX, or TXID")
	ErrInvalidLEI                               = errors.New("national identifier required with max length 35")
	ErrCompleteNationalIdentifierCountry        = errors.New("a legal person must not have a value for country if identifier type is not LEIX")
	ErrCompleteNationalIdentifierAuthorityEmpty = errors.New("a legal person must have a value for registration authority if identifier type is not LEIX")
	ErrCompleteNationalIdentifierAuthority      = errors.New("a legal person must not have a value for registration authority if identifier type is LEIX")
	ErrInvalidDateOfBirth                       = errors.New("date of birth must be a valid date in YYYY-MM-DD format")
	ErrInvalidPlaceOfBirth                      = errors.New("place of birth required with at most 70 characters")
	ErrDateInPast                               = errors.New("date of birth must be a historic date, prior to current date")
	ErrValidAddress                             = errors.New("address must have at least one address line or street name + building name or number")
	ErrInvalidAddressLines                      = errors.New("an address can contain at most 7 address lines")
)

// Parsing and JSON Serialization Errors
var (
	ErrPersonOneOfViolation              = errors.New("ivms101: person must be either a legal person or a natural person not both")
	ErrInvalidNaturalPersonNameTypeCode  = errors.New("ivms101: invalid natural person name type code")
	ErrParseNaturalPersonNameTypeCode    = errors.New("ivms101: could not parse natural person name type code from value")
	ErrInvalidLegalPersonNameTypeCode    = errors.New("ivms101: invalid legal person name type code")
	ErrParseLegalPersonNameTypeCode      = errors.New("ivms101: could not parse legal person name type code from value")
	ErrInvalidNationalIdentifierTypeCode = errors.New("ivms101: invalid national identifier type code")
	ErrParseNationalIdentifierTypeCode   = errors.New("ivms101: could not parse national identifier type code from value")
	ErrInvalidAddressTypeCode            = errors.New("ivms101: invalid address type code")
	ErrParseAddressTypeCode              = errors.New("ivms101: could not parse address type code from value")
	ErrInvalidTransliterationMethodCode  = errors.New("ivms101: invalid transliteration method code")
	ErrParseTransliterationMethodCode    = errors.New("ivms101: could not parse transliteration method code from value")
)

//===========================================================================
// Validation Errors
//===========================================================================

func MissingField(field string) *FieldError {
	return &FieldError{verb: "missing", field: field, issue: "this field is required"}
}

func IncorrectField(field, issue string) *FieldError {
	return &FieldError{verb: "invalid field", field: field, issue: issue}
}

func MaxNText(field string, max, length int) *FieldError {
	return &FieldError{field: fmt.Sprintf("exceeded max length %d chars", max), verb: field, issue: fmt.Sprintf("%d characters is too long", length)}
}

func InvalidEnum(field, enumValue, enumType string) *FieldError {
	return &FieldError{field: field, verb: "invalid enum", issue: fmt.Sprintf("%q is not a valid %s", enumValue, enumType)}
}

func OneOfMissing(fields ...string) *FieldError {
	var fieldstr string
	switch len(fields) {
	case 0:
		panic("no fields specified for one of")
	case 1:
		return MissingField(fields[0])
	default:
		fieldstr = fieldList(fields...)
	}

	return &FieldError{verb: "missing one of", field: fieldstr, issue: "at most one of these fields is required"}
}

func OneOfTooMany(fields ...string) *FieldError {
	if len(fields) < 2 {
		panic("must specify at least two fields for one of too many")
	}
	return &FieldError{verb: "specify only one of", field: fieldList(fields...), issue: "at most one of these fields may be specified"}
}

func ValidationError(prefix string, err error, errs ...error) error {
	var verr ValidationErrors
	if err == nil {
		verr = make(ValidationErrors, 0, len(errs))
	} else {
		switch v := err.(type) {
		case ValidationErrors:
			verr = v
		case *FieldError:
			verr = ValidationErrors{v}
		default:
			verr = make(ValidationErrors, 0, len(errs)+1)
			verr = append(verr, &FieldError{verb: "invalid", field: "input", issue: v.Error()})
		}
	}

	for _, e := range errs {
		if e != nil {
			switch v := e.(type) {
			case ValidationErrors:
				for _, serr := range v {
					verr = append(verr, serr.add(prefix))
				}
			case *FieldError:
				verr = append(verr, v.add(prefix))
			default:
				ierr := &FieldError{verb: "invalid", field: "input", issue: e.Error()}
				verr = append(verr, ierr.add(prefix))
			}
		}
	}

	if len(verr) == 0 {
		return nil
	}
	return verr
}

type ValidationErrors []*FieldError

func (e ValidationErrors) Error() string {
	if len(e) == 1 {
		return e[0].Error()
	}

	errs := make([]string, 0, len(e))
	for _, err := range e {
		errs = append(errs, err.Error())
	}

	return fmt.Sprintf("%d validation errors occurred:\n  %s", len(e), strings.Join(errs, "\n  "))
}

type FieldError struct {
	parent []string
	verb   string
	field  string
	issue  string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("ivms101: %s %s: %s", e.verb, e.Field(), e.issue)
}

func (e *FieldError) Field() string {
	if len(e.parent) > 0 {
		prefix := strings.Join(e.parent, ".")
		return fmt.Sprintf("%s.%s", prefix, e.field)
	}
	return e.field
}

func (e *FieldError) add(parent string) *FieldError {
	if parent != "" {
		if e.parent == nil {
			e.parent = make([]string, 0, 1)
		}
		e.parent = append([]string{parent}, e.parent...)
	}
	return e
}

func fieldList(fields ...string) string {
	switch len(fields) {
	case 0:
		return ""
	case 1:
		return fields[0]
	case 2:
		return fmt.Sprintf("%s or %s", fields[0], fields[1])
	default:
		last := len(fields) - 1
		return fmt.Sprintf("%s, or %s", strings.Join(fields[0:last], ", "), fields[last])
	}
}

//===========================================================================
// Wrapped Error
//===========================================================================

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
