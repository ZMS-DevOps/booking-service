package application

import (
	"encoding/json"
	"errors"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type ReservationRequestService struct {
	store                 domain.ReservationRequestStore
	unavailabilityService UnavailabilityService
	producer              *kafka.Producer
}

func NewReservationRequestService(store domain.ReservationRequestStore, unavailabilityService *UnavailabilityService, producer *kafka.Producer) *ReservationRequestService {
	return &ReservationRequestService{
		store:                 store,
		unavailabilityService: *unavailabilityService,
		producer:              producer,
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
		log.Printf("prosao2")
		service.produceNotification("reservation-request.created", reservationRequest.HostId, reservationRequest.Id.Hex(), "automatic")
	} else {
		service.produceNotification("reservation-request.created", reservationRequest.HostId, reservationRequest.Id.Hex(), "")
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
	log.Printf("stigao ovde")
	service.produceNotification("host-reviewed-reservation-request", reservationRequest.UserId, reservationRequest.Id.Hex(), "accept-request")
	log.Printf("prosao")
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

	reservationRequest.Status = domain.DeclinedByHost
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}
	service.produceNotification("host-reviewed-reservation-request", reservationRequest.UserId, reservationRequest.Id.Hex(), "decline-request")
	return nil
}

func (service *ReservationRequestService) DeleteRequest(id primitive.ObjectID) error {
	err := service.store.Delete(id)
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

	if !isReservationInFuture(reservationRequest) {
		return errors.New("reservation is in the future")
	}

	reservationRequest.Status = domain.DeclinedByUser
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}

	unavailabilityPeriod := domain.UnavailabilityPeriod{
		Start: reservationRequest.Start,
		End:   reservationRequest.End,
	}

	log.Printf("stigao DeclineReservation")
	log.Printf("unavailabilityPeriod start %s", unavailabilityPeriod.Start)
	log.Printf("unavailabilityPeriod end %s", unavailabilityPeriod.End)
	err = service.unavailabilityService.RemoveUnavailabilityPeriod(reservationRequest.AccommodationId, &unavailabilityPeriod, false)
	if err != nil {
		return err
	}

	log.Printf("stigao do produce")
	service.produceNotification("reservation.canceled", reservationRequest.HostId, reservationRequest.Id.Hex(), "canceled")

	log.Printf("prosao DeclineReservation")
	return nil
}

func (service *ReservationRequestService) DeleteClient(clientId string) bool {
	reservationRequests, err := service.store.GetByClientId(clientId)
	if err != nil {
		return false
	}
	for _, reservationRequest := range reservationRequests {
		if isReservationInFuture(reservationRequest) && reservationRequest.Status == domain.Approved {
			return false
		}
	}
	for _, reservationRequest := range reservationRequests {
		err = service.store.Delete(reservationRequest.Id)
		if err != nil {
			return false
		}
	}
	return true
}

func (service *ReservationRequestService) GetByClientId(clientId string, status *domain.ReservationRequestStatus) ([]*domain.ReservationRequest, error) {
	if status != nil {
		return service.store.GetByClientIdAndStatus(clientId, *status)
	} else {
		return service.store.GetByClientId(clientId)
	}

}

func (service *ReservationRequestService) GetNumberOfCanceled(clientId string) int {
	declinedRequests, err := service.store.GetByClientIdAndStatus(clientId, domain.DeclinedByUser)
	if err != nil {
		return 0
	}
	return len(declinedRequests)
}

func (service *ReservationRequestService) GetFilteredRequests(userId string, userType string, past bool, search string) ([]*domain.ReservationRequest, error) {
	var requests []*domain.ReservationRequest
	var err error

	if userType == "host" {
		requests, err = service.store.GetByHostAndTimeAndSearch(userId, past, search)
	} else {
		requests, err = service.store.GetByClientIdAndTimeAndSearch(userId, past, search)
	}

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func isReservationInFuture(reservationRequest *domain.ReservationRequest) bool {
	today := time.Now()
	if !today.Before(reservationRequest.Start.AddDate(0, 0, -1)) {
		return false
	}
	return true
}

func (service *ReservationRequestService) produceNotification(topic string, receiverId string, reservationId string, status string) {
	notificationDTO := dto.NotificationDTO{
		UserId:        receiverId,
		ReservationId: reservationId,
		Status:        status,
	}
	message, _ := json.Marshal(notificationDTO)
	err := service.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	}

	service.producer.Flush(4 * 1000)
}

func (service *ReservationRequestService) CheckGuestHasReservationForHost(reviewerId string, hostId string) bool {
	requests, err := service.store.GetPastAcceptedReservationRequestByClientIdAndHostId(reviewerId, hostId)
	if err != nil {
		return false
	}
	return requests != nil && len(requests) > 0
}

func (service *ReservationRequestService) CheckGuestHasReservationForAccommodation(reviewerId string, accommodationId primitive.ObjectID) bool {
	requests, err := service.store.GetByClientIdAndAccommodationId(reviewerId, accommodationId)
	if err != nil {
		return false
	}
	return requests != nil && len(requests) > 0
}

func (service *ReservationRequestService) CheckAccommodationHasReservation(accommodationId primitive.ObjectID) bool {
	status := domain.Approved
	reservationRequests, err := service.GetByAccommodationId(accommodationId, &status)
	if err != nil {
		return false
	}
	for _, reservationRequest := range reservationRequests {
		if isReservationInFuture(reservationRequest) {
			return false
		}
	}

	for _, reservationRequest := range reservationRequests {
		if err := service.DeleteRequest(reservationRequest.Id); err != nil {
			return false
		}
	}

	if err := service.unavailabilityService.DeleteByAccommodationId(accommodationId); err != nil {
		return false
	}

	return true
}
