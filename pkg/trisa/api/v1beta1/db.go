package api

import (
	"database/sql/driver"
	"fmt"

	"google.golang.org/protobuf/proto"
)

func (s *SecureEnvelope) Scan(src interface{}) error {
	// Convert src into a byte array for unmarshaling
	var source []byte
	switch t := src.(type) {
	case []byte:
		source = t
	case nil:
		return nil
	default:
		return fmt.Errorf("incompatible type for secure envelope: %T", t)
	}

	// Unmarshal the protocol buffers
	p := SecureEnvelope{}
	if err := proto.Unmarshal(source, &p); err != nil {
		return err
	}

	*s = p
	return nil
}

func (s *SecureEnvelope) Value() (_ driver.Value, err error) {
	if s == nil {
		return nil, nil
	}

	var data []byte
	if data, err = proto.Marshal(s); err != nil {
		return nil, err
	}

	return driver.Value(data), nil
}
