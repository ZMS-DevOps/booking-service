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

func MapRegistrationRequest(dto AddReservationRequestDto) *domain.ReservationRequest {
	return &domain.ReservationRequest{
		AccommodationId: dto.AccommodationId,
		UserId:          dto.UserId,
		Start:           dto.Start,
		End:             dto.End,
		NumberOfGuests:  dto.NumberOfGuests,
		PriceTotal:      dto.PriceTotal,
	}
}

func MapReservationRequestResponse(requests []*domain.ReservationRequest) []*ReservationRequestResponse {
	var response []*ReservationRequestResponse

	for i := range requests {
		reservationRequest := ReservationRequestResponse{
			Id:              requests[i].Id,
			AccommodationId: requests[i].AccommodationId,
			UserId:          requests[i].UserId,
			Start:           requests[i].Start,
			End:             requests[i].End,
			NumberOfGuests:  requests[i].NumberOfGuests,
			PriceTotal:      requests[i].PriceTotal,
			Status:          requests[i].Status,
		}
		response = append(response, &reservationRequest)
	}
	return response
}
