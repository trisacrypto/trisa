package models

import (
	"encoding/json"
	"errors"
	"strings"
)

func init() {
	for i, val := range VASPCategories {
		vaspcategories_name[i] = strings.ToUpper(val)
	}
}

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

const (
	VerificationNone          = VerificationState_NO_VERIFICATION
	VerificationSubmitted     = VerificationState_SUBMITTED
	VerificationEmailVerified = VerificationState_EMAIL_VERIFIED
	VerificationPending       = VerificationState_PENDING_REVIEW
	VerificationReviewed      = VerificationState_REVIEWED
	VerificationIssuing       = VerificationState_ISSUING_CERTIFICATE
	VerificationVerified      = VerificationState_VERIFIED
	VerificationRejected      = VerificationState_REJECTED
	VerificationAppealed      = VerificationState_APPEALED
	VerificationErrored       = VerificationState_ERRORED
)

const (
	ServiceStatusUnknown     = ServiceState_UNKNOWN
	ServiceStatusHealthy     = ServiceState_HEALTHY
	ServiceStatusUnhealthy   = ServiceState_UNHEALTHY
	ServiceStatusDanger      = ServiceState_DANGER
	ServiceStatusOffline     = ServiceState_OFFLINE
	ServiceStatusMaintenance = ServiceState_MAINTENANCE
)

var (
	vaspcategories_name = [14]string{}
	VASPCategories      = [14]string{
		VASPCategoryUnknown,
		VASPCategoryExchange,
		VASPCategoryDEX,
		VASPCategoryP2P,
		VASPCategoryKiosk,
		VASPCategoryCustodian,
		VASPCategoryOTC,
		VASPCategoryFund,
		VASPCategoryProject,
		VASPCategoryGambling,
		VASPCategoryMiner,
		VASPCategoryMixer,
		VASPCategoryIndividual,
		VASPCategoryOther,
	}
)

var (
	ErrUnknownBusinessCategory  = errors.New("could not parse business category from input")
	ErrUnknownVASPCategory      = errors.New("could not validate vasp category from input")
	ErrUnknownVerificationState = errors.New("could not parse verification state from input")
	ErrUnknownServiceState      = errors.New("could not parse service state from input")
)

func (b BusinessCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BusinessCategory) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return err
	}

	if *b, err = ParseBusinessCategory(val); err != nil {
		return err
	}
	return nil
}

// ParseBusinessCategory from text representation.
func ParseBusinessCategory(in any) (BusinessCategory, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(strings.ReplaceAll(val, " ", "_"))
		if code, ok := BusinessCategory_value[val]; ok {
			return BusinessCategory(code), nil
		}
	case int32:
		if _, ok := BusinessCategory_name[val]; ok {
			return BusinessCategory(val), nil
		}
	case float64:
		if _, ok := BusinessCategory_name[int32(val)]; ok {
			return BusinessCategory(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, err
		}
		if _, ok := BusinessCategory_name[int32(i)]; ok {
			return BusinessCategory(i), nil
		}
	}

	return 0, ErrUnknownBusinessCategory
}

// Validates a VASP category for a TRISA registration, if the input is not a valid
// VASP category an error is returned; otherwise the normalized category is returned.
func ValidVASPCategory(in string) (string, error) {
	in = strings.ToUpper(strings.ReplaceAll(in, " ", "_"))
	for i, cat := range vaspcategories_name {
		if in == cat {
			return VASPCategories[i], nil
		}
	}
	return "", ErrUnknownVASPCategory
}

func (b VerificationState) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *VerificationState) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return err
	}

	if *b, err = ParseVerificationState(val); err != nil {
		return err
	}
	return nil
}

func ParseVerificationState(in any) (VerificationState, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(strings.ReplaceAll(val, " ", "_"))
		if i, ok := VerificationState_value[val]; ok {
			return VerificationState(i), nil
		}
	case int32:
		if _, ok := VerificationState_name[val]; ok {
			return VerificationState(val), nil
		}
	case float64:
		if _, ok := VerificationState_name[int32(val)]; ok {
			return VerificationState(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, err
		}
		if _, ok := VerificationState_name[int32(i)]; ok {
			return VerificationState(i), nil
		}
	}
	return 0, ErrUnknownVerificationState
}

func (b ServiceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *ServiceState) UnmarshalJSON(data []byte) (err error) {
	var val interface{}
	if err = json.Unmarshal(data, &val); err != nil {
		return err
	}

	if *b, err = ParseServiceState(val); err != nil {
		return err
	}
	return nil
}

func ParseServiceState(in any) (ServiceState, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if i, ok := ServiceState_value[val]; ok {
			return ServiceState(i), nil
		}
	case int32:
		if _, ok := ServiceState_name[val]; ok {
			return ServiceState(val), nil
		}
	case float64:
		if _, ok := ServiceState_name[int32(val)]; ok {
			return ServiceState(val), nil
		}
	case json.Number:
		i, err := val.Int64()
		if err != nil {
			return 0, err
		}
		if _, ok := ServiceState_name[int32(i)]; ok {
			return ServiceState(i), nil
		}
	}

	return 0, ErrUnknownServiceState
}
