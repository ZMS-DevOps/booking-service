package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Unavailability struct {
	Id                                    primitive.ObjectID     `bson:"_id"`
	AccommodationId                       primitive.ObjectID     `bson:"accommodation_id"`
	AccommodationName                     string                 `bson:"accommodation_name"`
	HostId                                string                 `bson:"host_id"`
	UnavailabilityPeriods                 []UnavailabilityPeriod `bson:"unavailability_periods"`
	ReviewReservationRequestAutomatically bool                   `bson:"review_reservation_request_automatically"`
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
	Id                primitive.ObjectID       `bson:"_id"`
	AccommodationId   primitive.ObjectID       `bson:"accommodation_id"`
	AccommodationName string                   `bson:"accommodation_name"`
	HostId            string                   `bson:"host_id"`
	UserId            string                   `bson:"user_id"`
	Start             time.Time                `bson:"start"`
	End               time.Time                `bson:"end"`
	NumberOfGuests    int                      `bson:"number_of_guests"`
	PriceTotal        float32                  `bson:"price_total"`
	Status            ReservationRequestStatus `bson:"status"`
}

type ReservationRequestStatus int

const (
	Pending ReservationRequestStatus = iota
	Approved
	DeclinedByUser
	DeclinedByHost
	Completed
)
