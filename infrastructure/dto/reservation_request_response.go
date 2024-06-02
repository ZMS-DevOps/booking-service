package dto

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReservationRequestResponse struct {
	Id              primitive.ObjectID              `bson:"_id"`
	AccommodationId primitive.ObjectID              `bson:"accommodation_id"`
	UserId          primitive.ObjectID              `bson:"user_id"`
	Start           time.Time                       `bson:"start"`
	End             time.Time                       `bson:"end"`
	NumberOfGuests  int                             `bson:"number_of_guests"`
	PriceTotal      float32                         `bson:"price_total"`
	Status          domain.ReservationRequestStatus `bson:"status"`
}
