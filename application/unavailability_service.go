package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type UnavailabilityService struct {
	store                   domain.UnavailabilityStore
	reservationRequestStore domain.ReservationRequestStore
	producer                *kafka.Producer
}

func NewUnavailabilityService(store domain.UnavailabilityStore, producer *kafka.Producer, reservationRequestStore domain.ReservationRequestStore) *UnavailabilityService {
	return &UnavailabilityService{
		store:                   store,
		producer:                producer,
		reservationRequestStore: reservationRequestStore,
	}
}

func (service *UnavailabilityService) AddUnavailability(accommodationId primitive.ObjectID, accommodationName string, automatically bool, hostId string) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}

	if unavailability != nil {
		return fmt.Errorf("unavailability already exists for unavailability ID: %s", accommodationId.Hex())
	}

	newUnavailability := &domain.Unavailability{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       accommodationId,
		AccommodationName:                     accommodationName,
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		HostId:                                hostId,
		ReviewReservationRequestAutomatically: automatically,
	}

	if err := service.store.Insert(newUnavailability); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) UpdateUnavailability(accommodationId primitive.ObjectID, accommodationName string, automatically bool, hostId string) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}

	unavailability.ReviewReservationRequestAutomatically = automatically
	unavailability.AccommodationName = accommodationName
	unavailability.HostId = hostId

	if err := service.store.Update(unavailability.Id, unavailability); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) AddUnavailabilityPeriod(accommodationId primitive.ObjectID, period *domain.UnavailabilityPeriod) error {
	var unavailability, err = service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}
	period.Id = primitive.NewObjectID()

	if !(service.checkIfCouldAddUnavailability(unavailability, period)) {
		return errors.New("could not add unavailability period")
	}

	if err := service.store.UpdateUnavailabilityPeriods(unavailability.Id, insertPeriod(period, unavailability.UnavailabilityPeriods)); err != nil {
		return err
	}

	if err := service.reservationRequestStore.DeleteByAccommodation(accommodationId); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) checkIfCouldAddUnavailability(unavailability *domain.Unavailability, period *domain.UnavailabilityPeriod) bool {
	for _, unavailabilityPeriod := range unavailability.UnavailabilityPeriods {
		if periodsOverlap(unavailabilityPeriod.Start, unavailabilityPeriod.End, period.Start, period.End) {
			if unavailabilityPeriod.Reason == domain.Reserved {
				return false
			}

		}
	}
	return true
}

func (service *UnavailabilityService) RemoveUnavailabilityPeriod(accommodationId primitive.ObjectID, period *domain.UnavailabilityPeriod, shouldRemainReserved bool) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	log.Printf("unavailability %v removed from unavailability period %v\n", unavailability, period)
	if err != nil {
		return err
	}

	updatedPeriods := removePeriod(*period, unavailability.UnavailabilityPeriods, shouldRemainReserved)
	log.Printf("unavailability period %v removed from unavailability periods %v\n", updatedPeriods, period)
	if err := service.store.UpdateUnavailabilityPeriods(unavailability.Id, updatedPeriods); err != nil {
		return err
	}
	log.Printf("vrating period removed from unavailability periods %v\n", updatedPeriods)

	return nil
}

func (service *UnavailabilityService) GetAll() ([]*domain.Unavailability, error) {
	return service.store.GetAll()
}

func (service *UnavailabilityService) Get(id primitive.ObjectID) (*domain.Unavailability, error) {
	return service.store.Get(id)
}

func (service *UnavailabilityService) GetByAccommodationId(id primitive.ObjectID) (*domain.Unavailability, error) {
	return service.store.GetByAccommodationId(id)
}

func (service *UnavailabilityService) DeleteByAccommodationId(id primitive.ObjectID) error {
	return service.store.DeleteByAccommodationId(id)
}

func (service *UnavailabilityService) GetByHostId(id string) ([]*domain.Unavailability, error) {
	return service.store.GetByHostId(id)
}

func (service *UnavailabilityService) FilterAvailable(ids []primitive.ObjectID, startDate time.Time, endDate time.Time) ([]primitive.ObjectID, error) {
	var response []primitive.ObjectID
	for _, id := range ids {
		unavailability, err := service.store.GetByAccommodationId(id)
		if err != nil {
			return nil, err
		}
		if unavailability == nil {
			response = append(response, id)
			continue
		}
		available := true
		for _, period := range unavailability.UnavailabilityPeriods {
			if periodsOverlap(startDate, endDate, period.Start, period.End) {
				available = false
				break
			}
		}

		if available {
			response = append(response, id)
		}
	}
	return response, nil
}

func (service *UnavailabilityService) DeleteHost(hostId string) (bool, error) {
	unavailabilityList, err := service.store.GetByHostId(hostId)
	if err != nil {
		return false, err
	}
	for _, unavailability := range unavailabilityList {
		for _, period := range unavailability.UnavailabilityPeriods {
			log.Printf("start period %s\n", period.Start)
			log.Printf("reason %d\n", period.Reason)
			if isFuturePeriod(period) && period.Reason == domain.Reserved {
				return false, nil
			}
		}
	}
	log.Printf("dosao 1 deleteHost")

	service.produceDeleteAccommodationNotification(hostId)
	log.Printf("dosao 2 deleteHost")
	err = service.reservationRequestStore.DeleteByHost(hostId)
	log.Printf("dosao 3 deleteHost")
	if err != nil {
		return false, err
	}
	return true, nil
}

func (service *UnavailabilityService) IsAutomatic(accommodationId primitive.ObjectID) (bool, error) {
	var accommodation, err = service.store.GetByAccommodationId(accommodationId)
	log.Printf("accommodation: %v", accommodation)
	if err != nil {
		return false, err
	}
	return accommodation.ReviewReservationRequestAutomatically, nil
}

func isFuturePeriod(period domain.UnavailabilityPeriod) bool {
	now := time.Now()
	return period.Start.After(now)
}

func periodsOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}

func (service *UnavailabilityService) produceDeleteAccommodationNotification(hostId string) {
	var topic = "accommodation.delete"

	notificationDTO := dto.AccommodationDeleteNotification{
		Id: hostId,
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
