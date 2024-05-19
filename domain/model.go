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
