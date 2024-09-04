package ivms101

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func (i *IdentityPayload) Scan(src interface{}) error {
	if err := ScanJSON(src, i); err != nil {
		return err
	}
	return nil
}

func (i *IdentityPayload) Value() (_ driver.Value, err error) {
	return ValueJSON(i)
}

func (p *Person) Scan(src interface{}) error {
	if err := ScanJSON(src, p); err != nil {
		return err
	}
	return nil
}

func (p *Person) Value() (_ driver.Value, err error) {
	return ValueJSON(p)
}

func (p *NaturalPerson) Scan(src interface{}) error {
	if err := ScanJSON(src, p); err != nil {
		return err
	}
	return nil
}

func (p *NaturalPerson) Value() (_ driver.Value, err error) {
	return ValueJSON(p)
}

func (p *LegalPerson) Scan(src interface{}) error {
	if err := ScanJSON(src, p); err != nil {
		return err
	}
	return nil
}

func (p *LegalPerson) Value() (_ driver.Value, err error) {
	return ValueJSON(p)
}

func (a *Address) Scan(src interface{}) error {
	if err := ScanJSON(src, a); err != nil {
		return err
	}
	return nil
}

func (a *Address) Value() (_ driver.Value, err error) {
	return ValueJSON(a)
}

func ScanJSON(src, dst interface{}) error {
	// Convert src into a byte array to unmarshal json data
	var source []byte
	switch t := src.(type) {
	case []byte:
		source = t
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type to unmarshal json: %T", t)
	}

	if err := json.Unmarshal(source, dst); err != nil {
		return err
	}
	return nil
}

func ValueJSON(obj interface{}) (_ driver.Value, err error) {
	// Store null if obj is is nil
	if obj == nil {
		return nil, nil
	}

	// Store JSON as a BLOB for this type
	var data []byte
	if data, err = json.Marshal(obj); err != nil {
		return nil, err
	}
	return driver.Value(data), nil
}
