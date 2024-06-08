package application

import (
	"errors"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

	isAutomatic, err := service.unavailabilityService.IsAutomatic(reservationRequest.AccommodationId)
	if err != nil {
		return err
	}

	requestId, err := service.store.Insert(reservationRequest)

	if err != nil {
		return err
	}

	if isAutomatic {
		err = service.ApproveRequest(*requestId)
	}

	if err != nil {
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
	err = service.store.CancelOverlappingPendingRequests(reservationRequest)
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

func (service *ReservationRequestService) DeclineReservation(id primitive.ObjectID) error {
	reservationRequest, err := service.store.Get(id)
	if err != nil {
		return err
	}
	if reservationRequest.Status != domain.Approved {
		return errors.New("reservation is not approved")
	}

	err = isBeforeReservation(reservationRequest)
	if err != nil {
		return err
	}

	reservationRequest.Status = domain.Declined
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}

	unavailabilityPeriod := domain.UnavailabilityPeriod{
		Start: reservationRequest.Start,
		End:   reservationRequest.End,
	}

	err = service.unavailabilityService.RemoveUnavailabilityPeriod(reservationRequest.AccommodationId, &unavailabilityPeriod)
	if err != nil {
		return err
	}

	return nil
}

func (service *ReservationRequestService) DeleteClient(clientId primitive.ObjectID) bool {
	reservationRequests, err := service.store.GetByClientId(clientId)
	if err != nil {
		return false
	}
	for _, reservationRequest := range reservationRequests {
		if reservationRequest.Status == domain.Pending {
			return false
		}
	}
	for _, reservationRequest := range reservationRequests {
		if reservationRequest.Status == domain.Pending {
			err = service.store.Delete(reservationRequest.Id)
			if err != nil {
				return false
			}
		}
	}
	return true
}

func (service *ReservationRequestService) GetByClientId(clientId primitive.ObjectID, status *domain.ReservationRequestStatus) ([]*domain.ReservationRequest, error) {
	if status != nil {
		return service.store.GetByClientIdAndStatus(clientId, *status)
	} else {
		return service.store.GetByClientId(clientId)
	}

}

func (service *ReservationRequestService) GetNumberOfCanceled(clientId primitive.ObjectID) int {
	declinedRequests, err := service.store.GetByClientIdAndStatus(clientId, domain.Declined)
	if err != nil {
		return 0
	}
	return len(declinedRequests)
}

func isBeforeReservation(reservationRequest *domain.ReservationRequest) error {
	today := time.Now()
	if !today.Before(reservationRequest.Start.AddDate(0, 0, -1)) {
		return errors.New("cannot decline the reservation less than one day before the start date")
	}
	return nil
}
