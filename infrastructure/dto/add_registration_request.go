package dto

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AddReservationRequestDto struct {
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

func ValidateAddRegistrationRequestDto(dto AddReservationRequestDto) error {
	validate := validator.New()

	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Field '%s' failed validation with tag '%s'\n", err.Field(), err.Tag())
		}
		return err
	}

	return nil
}
