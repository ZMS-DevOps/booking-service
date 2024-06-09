package dto

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UnavailabilityResponse struct {
	Id                string `json:"id"`
	AccommodationId   string `json:"accommodation_id"`
	AccommodationName string `json:"accommodation_name"`
	Start             string `json:"start"`
	End               string `json:"end"`
	Reason            string `json:"reason"`
}

func MapToUnavailabilityResponse(id primitive.ObjectID, accommodationId primitive.ObjectID, accommodationName string, start time.Time, end time.Time, reason domain.UnavailabilityReason) UnavailabilityResponse {
	return UnavailabilityResponse{
		Id:                id.Hex(),
		AccommodationId:   accommodationId.Hex(),
		AccommodationName: accommodationName,
		Start:             start.Format(time.RFC3339),
		End:               end.Format(time.RFC3339),
		Reason:            unavailabilityReasonToString(reason),
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
