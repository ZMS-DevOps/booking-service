package dto

import (
	"github.com/ZMS-DevOps/booking-service/domain"
)

func MapUnavailabilityPeriod(addUnavailabilityPeriodDto *ManageUnavailabilityPeriodDto) *domain.UnavailabilityPeriod {
	unavailabilityPeriod := &domain.UnavailabilityPeriod{
		Start: addUnavailabilityPeriodDto.Start,
		End:   addUnavailabilityPeriodDto.End,
	}
	return unavailabilityPeriod
}

func MapRegistrationRequest(dto AddReservationRequestDto) *domain.ReservationRequest {
	return &domain.ReservationRequest{
		AccommodationId:   dto.AccommodationId,
		AccommodationName: dto.AccommodationName,
		UserId:            dto.UserId,
		HostId:            dto.HostId,
		Start:             dto.Start,
		End:               dto.End,
		NumberOfGuests:    dto.NumberOfGuests,
		PriceTotal:        dto.PriceTotal,
	}
}

func MapReservationRequestResponse(requests []*domain.ReservationRequest) []*ReservationRequestResponse {
	var response []*ReservationRequestResponse

	for i := range requests {
		reservationRequest := ReservationRequestResponse{
			Id:                requests[i].Id,
			AccommodationId:   requests[i].AccommodationId,
			AccommodationName: requests[i].AccommodationName,
			UserId:            requests[i].UserId,
			Start:             requests[i].Start,
			End:               requests[i].End,
			NumberOfGuests:    requests[i].NumberOfGuests,
			PriceTotal:        requests[i].PriceTotal,
			Status:            requests[i].Status,
		}
		response = append(response, &reservationRequest)
	}
	return response
}
