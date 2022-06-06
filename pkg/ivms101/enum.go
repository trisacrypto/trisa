package ivms101

import (
	"encoding/json"
	"errors"
	"strings"
)

//
// Constant Code Helpers
//

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
// NaturalPersonNameTypeCode JSON
//

const naturalPersonTypeCodePrefix = "NATURAL_PERSON_NAME_TYPE_CODE_"

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (n NaturalPersonNameTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(n.String(), naturalPersonTypeCodePrefix)
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (n *NaturalPersonNameTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse NaturalPersonNameTypeCode from value")
	}
	s = naturalPersonTypeCodePrefix + strings.ToUpper(s)
	code, ok := NaturalPersonNameTypeCode_value[s]
	if !ok {
		return errors.New("invalid NaturalPersonNameTypeCode alias")
	}
	*n = NaturalPersonNameTypeCode(code)
	return nil
}

//
// LegalPersonNameTypeCode JSON
//

const legalPersonNameTypeCodePrefix = "LEGAL_PERSON_NAME_TYPE_CODE_"

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (l LegalPersonNameTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(l.String(), legalPersonNameTypeCodePrefix)
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (l *LegalPersonNameTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse LegalPersonNameTypeCode from value")
	}
	s = legalPersonNameTypeCodePrefix + strings.ToUpper(s)
	code, ok := LegalPersonNameTypeCode_value[s]
	if !ok {
		return errors.New("invalid LegalPersonNameTypeCode alias")
	}
	*l = LegalPersonNameTypeCode(code)
	return nil
}

//
// AddressTypeCode JSON
//

const addressTypeCodePrefix = "ADDRESS_TYPE_CODE_"

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (a AddressTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(a.String(), addressTypeCodePrefix)
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (a *AddressTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse AddressTypeCode from value")
	}
	s = addressTypeCodePrefix + strings.ToUpper(s)
	code, ok := AddressTypeCode_value[s]
	if !ok {
		return errors.New("invalid AddressTypeCode alias")
	}
	*a = AddressTypeCode(code)
	return nil
}

//
// NationalIdentifierTypeCode JSON
//

const nationalIdentifierTypeCodePrefix = "NATIONAL_IDENTIFIER_TYPE_CODE_"

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (i NationalIdentifierTypeCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(i.String(), nationalIdentifierTypeCodePrefix)
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (i *NationalIdentifierTypeCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse NationalIdentifierTypeCode from value")
	}
	s = nationalIdentifierTypeCodePrefix + strings.ToUpper(s)
	code, ok := NationalIdentifierTypeCode_value[s]
	if !ok {
		return errors.New("invalid NationalIdentifierTypeCode alias")
	}
	*i = NationalIdentifierTypeCode(code)
	return nil
}

//
// TransliterationMethodCode JSON
//

const transliterationMethodCodePrefix = "TRANSLITERATION_METHOD_CODE_"

// Must be a value receiver to ensure it is marshaled correctly from it's parent struct
func (t TransliterationMethodCode) MarshalJSON() ([]byte, error) {
	data := strings.TrimPrefix(t.String(), transliterationMethodCodePrefix)
	return json.Marshal(data)
}

// Must be a pointer receiver so that we can indirect back to the correct variable
func (t *TransliterationMethodCode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return errors.New("could not parse TransliterationMethodCode from value")
	}
	s = transliterationMethodCodePrefix + strings.ToUpper(s)
	code, ok := TransliterationMethodCode_value[s]
	if !ok {
		return errors.New("invalid TransliterationMethodCode alias")
	}
	*t = TransliterationMethodCode(code)
	return nil
}
