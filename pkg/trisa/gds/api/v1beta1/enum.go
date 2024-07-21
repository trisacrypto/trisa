package api

import (
	"errors"
	"strings"
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
	ErrUnknownServiceState = errors.New("could not parse service state from input")
)

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
