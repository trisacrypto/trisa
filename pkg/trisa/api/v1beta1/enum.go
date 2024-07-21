package api

import (
	"errors"
	"strings"
)

//===========================================================================
// ENUM Constant Code Helpers
//===========================================================================

const (
	TransferStateUnspecified = TransferState_UNSPECIFIED
	TransferStarted          = TransferState_STARTED
	TransferPending          = TransferState_PENDING
	TransferRepair           = TransferState_REPAIR
	TransferReview           = TransferState_REVIEW
	TransferAccepted         = TransferState_ACCEPTED
	TransferCompleted        = TransferState_COMPLETED
	TransferRejected         = TransferState_REJECTED
)

const (
	ConfirmationTypeUnknown  = ConfirmationType_UNKNOWN
	ConfirmationTypeSimple   = ConfirmationType_SIMPLE
	ConfirmationTypeKeyToken = ConfirmationType_KEYTOKEN
	ConfirmationTypeOnChain  = ConfirmationType_ONCHAIN
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
	ErrUnknownTransferState    = errors.New("could not parse transfer state from input")
	ErrUnknownConfirmationType = errors.New("could not parse confirmation type from input")
	ErrUnknownServiceState     = errors.New("could not parse service state from input")
)

func ParseTransferState(in any) (TransferState, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if i, ok := TransferState_value[val]; ok {
			return TransferState(i), nil
		}
	case int32:
		if _, ok := TransferState_name[val]; ok {
			return TransferState(val), nil
		}
	}

	return 0, ErrUnknownTransferState
}

func ParseConfirmationType(in any) (ConfirmationType, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if i, ok := ConfirmationType_value[val]; ok {
			return ConfirmationType(i), nil
		}
	case int32:
		if _, ok := ConfirmationType_name[val]; ok {
			return ConfirmationType(val), nil
		}
	}

	return 0, ErrUnknownConfirmationType
}

func ParseServiceState(in any) (ServiceState_Status, error) {
	switch val := in.(type) {
	case string:
		val = strings.ToUpper(val)
		if i, ok := ServiceState_Status_value[val]; ok {
			return ServiceState_Status(i), nil
		}
	case int32:
		if _, ok := ServiceState_Status_name[val]; ok {
			return ServiceState_Status(val), nil
		}
	}

	return 0, ErrUnknownServiceState
}
