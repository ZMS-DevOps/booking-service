package dto

import (
	"fmt"
	"github.com/ZMS-DevOps/booking-service/domain"
)

func MapUnavailabilityPeriod(addUnavailabilityPeriodDto *ManageUnavailabilityPeriodDto) *domain.UnavailabilityPeriod {
	var reason, _ = mapReason(addUnavailabilityPeriodDto.Reason)
	unavailabilityPeriod := &domain.UnavailabilityPeriod{
		Start:  addUnavailabilityPeriodDto.Start,
		End:    addUnavailabilityPeriodDto.End,
		Reason: *reason,
	}
	return unavailabilityPeriod
}

func mapReason(reason *string) (*domain.UnavailabilityReason, error) {
	var unavailabilityReason domain.UnavailabilityReason
	switch *reason {
	case "Reserved":
		unavailabilityReason = domain.Reserved
	case "OwnerSet":
		unavailabilityReason = domain.OwnerSet
	default:
		return nil, fmt.Errorf("invalid pricing type: %s", *reason)
	}
	return &unavailabilityReason, nil
}
