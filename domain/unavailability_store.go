package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UnavailabilityStore interface {
	Get(id primitive.ObjectID) (*Unavailability, error)
	Insert(unavailability *Unavailability) error
	DeleteAll()
	GetAll() ([]*Unavailability, error)
	GetPeriod(id primitive.ObjectID) (UnavailabilityPeriod, error)

	GetUnavailabilityPeriods(id primitive.ObjectID) ([]UnavailabilityPeriod, error)
	UpdateUnavailabilityPeriods(unavailabilityId primitive.ObjectID, periods []UnavailabilityPeriod) error
	GetByAccommodationId(accommodationId primitive.ObjectID) (*Unavailability, error)
}
