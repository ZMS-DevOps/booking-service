package application

import (
	"fmt"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UnavailabilityService struct {
	store domain.UnavailabilityStore
}

func NewUnavailabilityService(store domain.UnavailabilityStore) *UnavailabilityService {
	return &UnavailabilityService{
		store: store,
	}
}

func (service *UnavailabilityService) AddUnavailability(accommodationId primitive.ObjectID, automatically bool) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}

	if unavailability != nil {
		fmt.Printf(unavailability.AccommodationId.Hex())
		return fmt.Errorf("unavailability already exists for unavailability ID: %s", accommodationId.Hex())
	}

	newUnavailability := &domain.Unavailability{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       accommodationId,
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: automatically,
	}

	if err := service.store.Insert(newUnavailability); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) UpdateUnavailability(accommodationId primitive.ObjectID, automatically bool) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}

	unavailability.ReviewReservationRequestAutomatically = automatically

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
	if err := service.store.UpdateUnavailabilityPeriods(unavailability.Id, insertPeriod(period, unavailability.UnavailabilityPeriods)); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) RemoveUnavailabilityPeriod(accommodationId primitive.ObjectID, period *domain.UnavailabilityPeriod) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}
	period.Id = primitive.NewObjectID()
	if err := service.store.UpdateUnavailabilityPeriods(unavailability.Id, removePeriod(*period, unavailability.UnavailabilityPeriods)); err != nil {
		return err
	}

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

func (service *UnavailabilityService) FilterAvailable(ids []primitive.ObjectID, startDate time.Time, endDate time.Time) ([]primitive.ObjectID, error) {
	var response []primitive.ObjectID
	for _, id := range ids {
		unavailability, err := service.store.GetByAccommodationId(id)
		if err != nil {
			return nil, err
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

func periodsOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.After(start2) && start1.Before(end2) ||
		start1.Before(start2) && end1.After(start2)
}
