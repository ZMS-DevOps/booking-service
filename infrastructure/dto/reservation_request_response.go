package dto

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ReservationRequestResponse struct {
	Id                           primitive.ObjectID              `json:"id"`
	AccommodationId              primitive.ObjectID              `json:"accommodation_id"`
	AccommodationName            string                          `json:"accommodation_name"`
	UserId                       primitive.ObjectID              `json:"user_id"`
	Start                        time.Time                       `json:"start"`
	End                          time.Time                       `json:"end"`
	NumberOfGuests               int                             `json:"number_of_guests"`
	PriceTotal                   float32                         `json:"price_total"`
	Status                       domain.ReservationRequestStatus `json:"status"`
	NumberOfCanceledReservations int                             `json:"number_of_canceled_reservations"`
}
