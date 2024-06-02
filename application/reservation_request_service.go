package application

import (
	"errors"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationRequestService struct {
	store                 domain.ReservationRequestStore
	unavailabilityService UnavailabilityService
}

func NewReservationRequestService(store domain.ReservationRequestStore, unavailabilityService *UnavailabilityService) *ReservationRequestService {
	return &ReservationRequestService{
		store:                 store,
		unavailabilityService: *unavailabilityService,
	}
}

func (service *ReservationRequestService) AddReservationRequest(reservationRequest *domain.ReservationRequest) error {
	reservationRequest.Id = primitive.NewObjectID()
	reservationRequest.Status = domain.Pending

	// todo add check if automatic and check if date range is available
	if err := service.store.Insert(reservationRequest); err != nil {
		return err
	}
	return nil
}

func (service *ReservationRequestService) GetByAccommodationId(accommodationId primitive.ObjectID, requestType *domain.ReservationRequestStatus) ([]*domain.ReservationRequest, error) {
	if requestType == nil {
		return service.store.GetByAccommodationId(accommodationId)
	} else {
		return service.store.GetByAccommodationIdAndType(accommodationId, *requestType)
	}
}

func (service *ReservationRequestService) ApproveRequest(id primitive.ObjectID) error {
	reservationRequest, err := service.store.Get(id)
	if err != nil {
		return err
	}
	if reservationRequest.Status != domain.Pending {
		return errors.New("reservation is not pending")
	}

	reservationRequest.Status = domain.Approved
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}
	// block all other pending requests that have overlapping periods
	err = service.store.CancelOverlappingPendingRequests(reservationRequest)
	// create period of unavailability
	err = service.createUnavailabilityPeriod(reservationRequest)
	if err != nil {
		return err
	}
	return nil
}

func (service *ReservationRequestService) DeclineRequest(id primitive.ObjectID) error {
	reservationRequest, err := service.store.Get(id)
	if err != nil {
		return err
	}
	if reservationRequest.Status != domain.Pending {
		return errors.New("reservation is not pending")
	}

	reservationRequest.Status = domain.Declined
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}
	return nil
}

func (service *ReservationRequestService) DeleteRequest(id primitive.ObjectID) error {
	_, err := service.store.Get(id)
	if err != nil {
		return err
	}
	err = service.store.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *ReservationRequestService) createUnavailabilityPeriod(reservationRequest *domain.ReservationRequest) error {
	unavailabilityPeriod := domain.UnavailabilityPeriod{
		Start:  reservationRequest.Start,
		End:    reservationRequest.End,
		Reason: domain.Reserved,
	}
	err := service.unavailabilityService.AddUnavailabilityPeriod(reservationRequest.AccommodationId, &unavailabilityPeriod)
	return err
}
