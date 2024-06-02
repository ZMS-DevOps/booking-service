package application

import (
	"fmt"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnavailabilityService struct {
	store domain.UnavailabilityStore
}

func NewUnavailabilityService(store domain.UnavailabilityStore) *UnavailabilityService {
	return &UnavailabilityService{
		store: store,
	}
}

func (service *UnavailabilityService) AddUnavailability(accommodationId primitive.ObjectID) error {
	unavailability, err := service.store.GetByAccommodationId(accommodationId)
	if err != nil {
		return err
	}

	if unavailability != nil {
		fmt.Printf(unavailability.AccommodationId.Hex())
		return fmt.Errorf("unavailability already exists for accommodation ID: %s", accommodationId.Hex())
	}

	newUnavailability := &domain.Unavailability{
		Id:                    primitive.NewObjectID(),
		AccommodationId:       accommodationId,
		UnavailabilityPeriods: []domain.UnavailabilityPeriod{},
	}

	if err := service.store.Insert(newUnavailability); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) AddUnavailabilityPeriod(unavailabilityId primitive.ObjectID, period *domain.UnavailabilityPeriod) error {
	_, err := service.store.Get(unavailabilityId)
	if err != nil {
		return err
	}

	var periods, _ = service.store.GetUnavailabilityPeriods(unavailabilityId)
	period.Id = primitive.NewObjectID()
	if err := service.store.UpdateUnavailabilityPeriods(unavailabilityId, insertPeriod(period, periods)); err != nil {
		return err
	}

	return nil
}

func (service *UnavailabilityService) RemoveUnavailabilityPeriod(unavailabilityId primitive.ObjectID, period *domain.UnavailabilityPeriod) interface{} {
	_, err := service.store.Get(unavailabilityId)
	if err != nil {
		return err
	}

	var periods, _ = service.store.GetUnavailabilityPeriods(unavailabilityId)
	period.Id = primitive.NewObjectID()
	if err := service.store.UpdateUnavailabilityPeriods(unavailabilityId, removePeriod(*period, periods)); err != nil {
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