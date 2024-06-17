package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UnavailabilityStore interface {
	Get(id primitive.ObjectID) (*Unavailability, error)
	Insert(unavailability *Unavailability) error
	DeleteAll()
	DeleteByAccommodationId(accommodationId primitive.ObjectID) error
	GetAll() ([]*Unavailability, error)
	GetPeriod(id primitive.ObjectID) (UnavailabilityPeriod, error)
	Update(id primitive.ObjectID, unavailability *Unavailability) error
	GetUnavailabilityPeriods(id primitive.ObjectID) ([]UnavailabilityPeriod, error)
	UpdateUnavailabilityPeriods(unavailabilityId primitive.ObjectID, periods []UnavailabilityPeriod) error
	GetByAccommodationId(accommodationId primitive.ObjectID) (*Unavailability, error)
	GetByHostId(id string) ([]*Unavailability, error)
}
