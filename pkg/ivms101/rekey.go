package ivms101

import (
	"encoding/json"
	"fmt"
)

var (
	allowRekeying         bool
	disallowUnknownFields bool
	nullJSON              = []byte{110, 117, 108, 108}
)

// Rekeying sets the module to perform a rekeying operation that changes snake_case and
// pluralized keys to the OpenVASP compatible keys in a JSON payload. This incurs some
// extra cost when unmarshaling but can make loading IVMS101 JSON data more flexible.
//
// Rekeying is false by default.
func AllowRekeying() {
	allowRekeying = true
}

// Turns of rekeying during JSON unmarshaling.
func DisallowRekeying() {
	allowRekeying = false
	disallowUnknownFields = false
}

// DisallowUnknownFields will return an error when an unknown field is part of the JSON
// message; this may make data validation easier.
//
// DisallowUnknownFields is false by default. Rekeying must be allowed to disallow
// unknown fields and rekeying is automatically set to true when this method is called.
func DisallowUnknownFields() {
	allowRekeying = true
	disallowUnknownFields = true
}

// Turns of disallow unknown fields checking during JSON unmarshaling.
func AllowUnknownFields() {
	disallowUnknownFields = false
}

func Rekey(data []byte, keyMap map[string]string) (_ []byte, err error) {
	if !allowRekeying {
		return data, nil
	}

	var message map[string]json.RawMessage
	if err = json.Unmarshal(data, &message); err != nil {
		return nil, err
	}

	for key, raw := range message {
		rekey, ok := keyMap[key]
		if !ok {
			if disallowUnknownFields {
				return nil, fmt.Errorf("json: unknown field %q", key)
			}
			continue
		}

		message[rekey] = raw
	}

	return json.Marshal(message)
}
