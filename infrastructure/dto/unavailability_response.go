package dto

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"time"
)

type UnavailabilityResponse struct {
	Id                                    string                    `json:"id"`
	AccommodationId                       string                    `json:"accommodation_id"`
	UnavailabilityPeriods                 []UnavailabilityPeriodDTO `json:"unavailabilityPeriods"`
	ReviewReservationRequestAutomatically bool                      `json:"reviewReservationRequestAutomatically"`
}

type UnavailabilityPeriodDTO struct {
	Id     string `json:"id"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Reason string `json:"reason"`
}

func MapToUnavailabilityResponse(unavailability domain.Unavailability) UnavailabilityResponse {
	periodsDTO := make([]UnavailabilityPeriodDTO, len(unavailability.UnavailabilityPeriods))
	for i, period := range unavailability.UnavailabilityPeriods {
		periodsDTO[i] = UnavailabilityPeriodDTO{
			Id:     period.Id.Hex(),
			Start:  period.Start.Format(time.RFC3339),
			End:    period.End.Format(time.RFC3339),
			Reason: unavailabilityReasonToString(period.Reason),
		}
	}

	return UnavailabilityResponse{
		Id:                                    unavailability.Id.Hex(),
		AccommodationId:                       unavailability.AccommodationId.Hex(),
		UnavailabilityPeriods:                 periodsDTO,
		ReviewReservationRequestAutomatically: unavailability.ReviewReservationRequestAutomatically,
	}
}

func unavailabilityReasonToString(reason domain.UnavailabilityReason) string {
	switch reason {
	case domain.Reserved:
		return "Reserved"
	case domain.OwnerSet:
		return "OwnerSet"
	default:
		return "Unknown"
	}
}
