package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Unavailability struct {
	Id                    primitive.ObjectID     `bson:"_id"`
	AccommodationId       primitive.ObjectID     `bson:"accommodation_id"`
	UnavailabilityPeriods []UnavailabilityPeriod `bson:"unavailability_periods"`
}

type UnavailabilityPeriod struct {
	Id     primitive.ObjectID   `bson:"_id"`
	Start  time.Time            `bson:"start"`
	End    time.Time            `bson:"end"`
	Reason UnavailabilityReason `bson:"reason"`
}

type UnavailabilityReason int

const (
	Reserved UnavailabilityReason = iota
	OwnerSet
)

type ReservationRequest struct {
	Id              primitive.ObjectID `bson:"_id"`
	AccommodationId primitive.ObjectID `bson:"accommodation_id"`
	UserId          primitive.ObjectID `bson:"user_id"`
	Start           time.Time          `bson:"start"`
	End             time.Time          `bson:"end"`
	NumberOfGuests  int                `bson:"number_of_guests"`
	PriceTotal      float32            `bson:"price_total"`
}

type ReservationRequestStatus int

const (
	Pending ReservationRequestStatus = iota
	Approved
	Declined
)