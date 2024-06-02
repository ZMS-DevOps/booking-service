package dto

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ManageUnavailabilityPeriodDto struct {
	UnavailabilityId primitive.ObjectID `json:"unavailability_id" validate:"required"`
	Start            time.Time          `json:"start" validate:"required"`
	End              time.Time          `json:"end" validate:"required,gtefield=Start"`
	Reason           *string            `json:"reason" validate:"omitempty,oneof=Reserved OwnerSet"`
}

func ValidateManageUnavailabilityPeriodDto(dto ManageUnavailabilityPeriodDto) error {
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