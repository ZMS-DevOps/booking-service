package dto

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AddReservationRequestDto struct {
	AccommodationId   primitive.ObjectID `json:"accommodation_id"`
	AccommodationName string             `json:"accommodation_name"`
	HostId            string             `json:"host_id"`
	UserId            string             `json:"user_id"`
	Start             time.Time          `json:"start"`
	End               time.Time          `json:"end"`
	NumberOfGuests    int                `json:"number_of_guests"`
	PriceTotal        float32            `json:"price_total"`
}

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
