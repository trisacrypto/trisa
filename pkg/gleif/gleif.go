package gleif

import (
	"errors"
	"regexp"
)

type RegistrationAuthority struct {
	Option       string `json:"option"`
	CountryName  string `json:"country_name"`
	Country      string `json:"country"`
	Jurisdiction string `json:"jurisdiction"`
	Register     string `json:"register"`
	Organization string `json:"organization"`
	Website      string `json:"website"`
	Comments     string `json:"comments"`
}

type RegistrationAuthorities []*RegistrationAuthority

func init() {
	validRAs = make(map[string]struct{}, len(registrationAuthorities))
	for _, ra := range registrationAuthorities {
		validRAs[ra.Option] = struct{}{}
	}
}

var (
	validRAs map[string]struct{}
	raregex  = regexp.MustCompile(`^RA[0-9]{6}$`)
)

var (
	ErrNotFound        = errors.New("gleif: registration authority not found")
	ErrIncorrectFormat = errors.New("gleif: invalid registration authority format")
)

func Find(ra, country string) (RegistrationAuthority, error) {
	if !raregex.MatchString(ra) {
		return RegistrationAuthority{}, ErrIncorrectFormat
	}

	if _, ok := validRAs[ra]; !ok {
		return RegistrationAuthority{}, ErrNotFound
	}

	for _, found := range registrationAuthorities {
		if found.Option == ra {
			if country == "" || found.Country == country {
				return *found, nil
			}
		}
	}

	return RegistrationAuthority{}, ErrNotFound
}

func Validate(ra string) error {
	_, err := Find(ra, "")
	return err
}
