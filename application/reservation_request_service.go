package application

import (
	"encoding/json"
	"errors"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/ZMS-DevOps/booking-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type ReservationRequestService struct {
	store                 domain.ReservationRequestStore
	unavailabilityService UnavailabilityService
	producer              *kafka.Producer
	loki                  promtail.Client
}

func NewReservationRequestService(store domain.ReservationRequestStore, unavailabilityService *UnavailabilityService, producer *kafka.Producer, loki promtail.Client) *ReservationRequestService {
	return &ReservationRequestService{
		store:                 store,
		unavailabilityService: *unavailabilityService,
		producer:              producer,
		loki:                  loki,
	}
}

func (service *ReservationRequestService) AddReservationRequest(reservationRequest *domain.ReservationRequest, span trace.Span, loki promtail.Client) error {
	reservationRequest.Id = primitive.NewObjectID()
	reservationRequest.Status = domain.Pending

	isAutomatic, err := service.unavailabilityService.IsAutomatic(reservationRequest.AccommodationId, span, loki)
	if err != nil {
		return err
	}

	util.HttpTraceInfo("Adding reservation request...", span, loki, "AddReservationRequest", "")
	requestId, err := service.store.Insert(reservationRequest)

	if err != nil {
		return err
	}

	if isAutomatic {
		err = service.ApproveRequest(*requestId, span, loki)
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

func (service *ReservationRequestService) GetByAccommodationId(accommodationId primitive.ObjectID, requestType *domain.ReservationRequestStatus, span trace.Span, loki promtail.Client) ([]*domain.ReservationRequest, error) {
	if requestType == nil {
		util.HttpTraceInfo("Fetching reservation requests by accommodation id...", span, loki, "GetByAccommodationId", "")
		return service.store.GetByAccommodationId(accommodationId)
	} else {
		util.HttpTraceInfo("Fetching reservation requests by accomodation id and type...", span, loki, "GetByAccommodationId", "")
		return service.store.GetByAccommodationIdAndType(accommodationId, *requestType)
	}
}

func (service *ReservationRequestService) ApproveRequest(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Fetching reservation requests by id...", span, loki, "GetByAccommodationId", "")
	reservationRequest, err := service.store.Get(id)
	if err != nil {
		return err
	}
	if reservationRequest.Status != domain.Pending {
		return errors.New("reservation is not pending")
	}

	reservationRequest.Status = domain.Approved
	util.HttpTraceInfo("Updating reservation requests...", span, loki, "GetByAccommodationId", "")
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}
	util.HttpTraceInfo("Canceling overlapping pending requests...", span, loki, "GetByAccommodationId", "")
	err = service.store.CancelOverlappingPendingRequests(reservationRequest)
	err = service.createUnavailabilityPeriod(reservationRequest, span, loki)
	if err != nil {
		return err
	}
	log.Printf("stigao ovde")
	service.produceNotification("host-reviewed-reservation-request", reservationRequest.UserId, reservationRequest.Id.Hex(), "accept-request")
	log.Printf("prosao")
	return nil
}

func (service *ReservationRequestService) DeclineRequest(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Fetching reservation requests by id...", span, loki, "DeclineRequest", "")
	reservationRequest, err := service.store.Get(id)
	if err != nil {
		return err
	}
	if reservationRequest.Status != domain.Pending {
		return errors.New("reservation is not pending")
	}

	reservationRequest.Status = domain.DeclinedByHost
	util.HttpTraceInfo("Updating reservation requests...", span, loki, "DeclineRequest", "")
	err = service.store.Update(id, reservationRequest)
	if err != nil {
		return err
	}
	service.produceNotification("host-reviewed-reservation-request", reservationRequest.UserId, reservationRequest.Id.Hex(), "decline-request")
	return nil
}

func (service *ReservationRequestService) DeleteRequest(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Deleting reservation requests...", span, loki, "DeleteRequest", "")
	err := service.store.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *ReservationRequestService) createUnavailabilityPeriod(reservationRequest *domain.ReservationRequest, span trace.Span, loki promtail.Client) error {
	unavailabilityPeriod := domain.UnavailabilityPeriod{
		Start:  reservationRequest.Start,
		End:    reservationRequest.End,
		Reason: domain.Reserved,
	}
	err := service.unavailabilityService.AddUnavailabilityPeriod(reservationRequest.AccommodationId, &unavailabilityPeriod, span, loki)
	return err
}

func (service *ReservationRequestService) DeclineReservation(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Fetching reservation requests by id...", span, loki, "DeclineReservation", "")
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
	util.HttpTraceInfo("Updating reservation requests...", span, loki, "DeclineReservation", "")
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
	err = service.unavailabilityService.RemoveUnavailabilityPeriod(reservationRequest.AccommodationId, &unavailabilityPeriod, false, span, loki)
	if err != nil {
		return err
	}

	log.Printf("stigao do produce")
	service.produceNotification("reservation.canceled", reservationRequest.HostId, reservationRequest.Id.Hex(), "canceled")

	log.Printf("prosao DeclineReservation")
	return nil
}

func (service *ReservationRequestService) DeleteClient(clientId string, span trace.Span, loki promtail.Client) bool {
	util.HttpTraceInfo("Fetching reservation requests by client id...", span, loki, "DeleteClient", "")
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
		util.HttpTraceInfo("Deleting reservation requests by id...", span, loki, "DeleteClient", "")
		err = service.store.Delete(reservationRequest.Id)
		if err != nil {
			return false
		}
	}
	return true
}

func (service *ReservationRequestService) GetByClientId(clientId string, status *domain.ReservationRequestStatus, span trace.Span, loki promtail.Client) ([]*domain.ReservationRequest, error) {
	if status != nil {
		util.HttpTraceInfo("Fetching reservation requests by client id and status...", span, loki, "GetByClientId", "")
		return service.store.GetByClientIdAndStatus(clientId, *status)
	} else {
		util.HttpTraceInfo("Fetching reservation requests by client id...", span, loki, "GetByClientId", "")
		return service.store.GetByClientId(clientId)
	}

}

func (service *ReservationRequestService) GetNumberOfCanceled(clientId string, span trace.Span, loki promtail.Client) int {
	util.HttpTraceInfo("Fetching accommodation by client id and status...", span, loki, "GetNumberOfCanceled", "")
	declinedRequests, err := service.store.GetByClientIdAndStatus(clientId, domain.DeclinedByUser)
	if err != nil {
		return 0
	}
	return len(declinedRequests)
}

func (service *ReservationRequestService) GetFilteredRequests(userId string, userType string, past bool, search string, span trace.Span, loki promtail.Client) ([]*domain.ReservationRequest, error) {
	var requests []*domain.ReservationRequest
	var err error

	if userType == "host" {
		util.HttpTraceInfo("Fetching reservation requests by host id and time...", span, loki, "GetFilteredRequests", "")
		requests, err = service.store.GetByHostAndTimeAndSearch(userId, past, search)
	} else {
		util.HttpTraceInfo("Fetching reservation requests by client id and time...", span, loki, "GetFilteredRequests", "")
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

func (service *ReservationRequestService) CheckGuestHasReservationForHost(reviewerId string, hostId string, span trace.Span, loki promtail.Client) bool {
	util.HttpTraceInfo("Fetching reservation requests by host id and accommodation id...", span, loki, "CheckGuestHasReservationForHost", "")
	requests, err := service.store.GetPastAcceptedReservationRequestByClientIdAndHostId(reviewerId, hostId)
	if err != nil {
		return false
	}
	return requests != nil && len(requests) > 0
}

func (service *ReservationRequestService) CheckGuestHasReservationForAccommodation(reviewerId string, accommodationId primitive.ObjectID, span trace.Span, loki promtail.Client) bool {
	util.HttpTraceInfo("Fetching reservation requests by client id and accommodation id...", span, loki, "CheckGuestHasReservationForAccommodation", "")
	requests, err := service.store.GetByClientIdAndAccommodationId(reviewerId, accommodationId)
	if err != nil {
		return false
	}
	return requests != nil && len(requests) > 0
}

func (service *ReservationRequestService) CheckAccommodationHasReservation(accommodationId primitive.ObjectID, span trace.Span, loki promtail.Client) bool {
	status := domain.Approved
	reservationRequests, err := service.GetByAccommodationId(accommodationId, &status, span, loki)
	if err != nil {
		return false
	}
	for _, reservationRequest := range reservationRequests {
		if isReservationInFuture(reservationRequest) {
			return false
		}
	}

	for _, reservationRequest := range reservationRequests {
		if err := service.DeleteRequest(reservationRequest.Id, span, loki); err != nil {
			return false
		}
	}

	if err := service.unavailabilityService.DeleteByAccommodationId(accommodationId, span, loki); err != nil {
		return false
	}

	return true
}
