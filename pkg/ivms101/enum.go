package ivms101

import (
	"encoding/json"
	"strings"
)

//===========================================================================
// ENUM Constant Code Helpers
//===========================================================================

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
	LegalPersonMisc    = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_MISC
	LegalPersonLegal   = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_LEGL
	LegalPersonShort   = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_SHRT
	LegalPersonTrading = LegalPersonNameTypeCode_LEGAL_PERSON_NAME_TYPE_CODE_TRAD
)

// Short form address type codes.
const (
	AddressTypeMisc       = AddressTypeCode_ADDRESS_TYPE_CODE_MISC
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

// Short form transliteration method codes.
const (
	TransliterationMethodOTHR = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_OTHR
	TransliterationMethodARAB = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_ARAB
	TransliterationMethodARAN = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_ARAN
	TransliterationMethodARMN = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_ARMN
	TransliterationMethodCYRL = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_CYRL
	TransliterationMethodDEVA = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_DEVA
	TransliterationMethodGEOR = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_GEOR
	TransliterationMethodGREK = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_GREK
	TransliterationMethodHANI = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_HANI
	TransliterationMethodHEBR = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_HEBR
	TransliterationMethodKANA = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_KANA
	TransliterationMethodKORE = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_KORE
	TransliterationMethodTHAI = TransliterationMethodCode_TRANSLITERATION_METHOD_CODE_THAI
)

//===========================================================================
// ENUM JSON Marshal and Unmarshal
//===========================================================================

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
func (n *NaturalPersonNameTypeCode) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return Wrap(ErrParseNaturalPersonNameTypeCode, err)
	}

	var code NaturalPersonNameTypeCode
	if code, err = ParseNaturalPersonNameTypeCode(val); err != nil {
		return err
	}

	*n = code
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
func (l *LegalPersonNameTypeCode) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return Wrap(ErrParseLegalPersonNameTypeCode, err)
	}

	var code LegalPersonNameTypeCode
	if code, err = ParseLegalPersonNameTypeCode(val); err != nil {
		return err
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
func (a *AddressTypeCode) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return Wrap(ErrParseAddressTypeCode, err)
	}

	var code AddressTypeCode
	if code, err = ParseAddressTypeCode(val); err != nil {
		return err
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
func (i *NationalIdentifierTypeCode) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return Wrap(ErrParseNationalIdentifierTypeCode, err)
	}

	var code NationalIdentifierTypeCode
	if code, err = ParseNationalIdentifierTypeCode(val); err != nil {
		return err
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
func (t *TransliterationMethodCode) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return Wrap(ErrParseTransliterationMethodCode, err)
	}

	var code TransliterationMethodCode
	if code, err = ParseTransliterationMethodCode(val); err != nil {
		return err
	}

	*t = TransliterationMethodCode(code)
	return nil
}

//===========================================================================
// Enum Parsing
//===========================================================================

func ParseNaturalPersonNameTypeCode(in any) (NaturalPersonNameTypeCode, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if !strings.HasPrefix(val, naturalPersonTypeCodePrefix) {
			val = naturalPersonTypeCodePrefix + val
		}

		if i, ok := NaturalPersonNameTypeCode_value[val]; ok {
			return NaturalPersonNameTypeCode(i), nil
		}
	case int32:
		if _, ok := NaturalPersonNameTypeCode_name[val]; ok {
			return NaturalPersonNameTypeCode(val), nil
		}
	case float64:
		if _, ok := NaturalPersonNameTypeCode_name[int32(val)]; ok {
			return NaturalPersonNameTypeCode(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, Wrap(ErrParseNaturalPersonNameTypeCode, err)
		}
		if _, ok := NaturalPersonNameTypeCode_name[int32(i)]; ok {
			return NaturalPersonNameTypeCode(i), nil
		}
	}

	return 0, ErrInvalidNaturalPersonNameTypeCode
}

func ParseLegalPersonNameTypeCode(in any) (LegalPersonNameTypeCode, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if !strings.HasPrefix(val, legalPersonNameTypeCodePrefix) {
			val = legalPersonNameTypeCodePrefix + val
		}

		if i, ok := LegalPersonNameTypeCode_value[val]; ok {
			return LegalPersonNameTypeCode(i), nil
		}
	case int32:
		if _, ok := LegalPersonNameTypeCode_name[val]; ok {
			return LegalPersonNameTypeCode(val), nil
		}
	case float64:
		if _, ok := LegalPersonNameTypeCode_name[int32(val)]; ok {
			return LegalPersonNameTypeCode(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, Wrap(ErrInvalidLegalPersonNameTypeCode, err)
		}
		if _, ok := LegalPersonNameTypeCode_name[int32(i)]; ok {
			return LegalPersonNameTypeCode(i), nil
		}
	}

	return 0, ErrInvalidLegalPersonNameTypeCode
}

func ParseAddressTypeCode(in any) (AddressTypeCode, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if !strings.HasPrefix(val, addressTypeCodePrefix) {
			val = addressTypeCodePrefix + val
		}

		if i, ok := AddressTypeCode_value[val]; ok {
			return AddressTypeCode(i), nil
		}
	case int32:
		if _, ok := AddressTypeCode_name[val]; ok {
			return AddressTypeCode(val), nil
		}
	case float64:
		if _, ok := AddressTypeCode_name[int32(val)]; ok {
			return AddressTypeCode(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, Wrap(ErrInvalidAddressTypeCode, err)
		}
		if _, ok := AddressTypeCode_name[int32(i)]; ok {
			return AddressTypeCode(i), nil
		}
	}

	return 0, ErrInvalidAddressTypeCode
}

func ParseNationalIdentifierTypeCode(in any) (NationalIdentifierTypeCode, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if !strings.HasPrefix(val, nationalIdentifierTypeCodePrefix) {
			val = nationalIdentifierTypeCodePrefix + val
		}

		if i, ok := NationalIdentifierTypeCode_value[val]; ok {
			return NationalIdentifierTypeCode(i), nil
		}
	case int32:
		if _, ok := NationalIdentifierTypeCode_name[val]; ok {
			return NationalIdentifierTypeCode(val), nil
		}
	case float64:
		if _, ok := NationalIdentifierTypeCode_name[int32(val)]; ok {
			return NationalIdentifierTypeCode(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, Wrap(ErrInvalidNationalIdentifierTypeCode, err)
		}
		if _, ok := NationalIdentifierTypeCode_name[int32(i)]; ok {
			return NationalIdentifierTypeCode(i), nil
		}
	}

	return 0, ErrInvalidNationalIdentifierTypeCode
}

func ParseTransliterationMethodCode(in any) (TransliterationMethodCode, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if !strings.HasPrefix(val, transliterationMethodCodePrefix) {
			val = transliterationMethodCodePrefix + val
		}

		if i, ok := TransliterationMethodCode_value[val]; ok {
			return TransliterationMethodCode(i), nil
		}
	case int32:
		if _, ok := TransliterationMethodCode_name[val]; ok {
			return TransliterationMethodCode(val), nil
		}
	case float64:
		if _, ok := TransliterationMethodCode_name[int32(val)]; ok {
			return TransliterationMethodCode(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, Wrap(ErrInvalidTransliterationMethodCode, err)
		}
		if _, ok := TransliterationMethodCode_name[int32(i)]; ok {
			return TransliterationMethodCode(i), nil
		}
	}

	return 0, ErrInvalidTransliterationMethodCode
}
