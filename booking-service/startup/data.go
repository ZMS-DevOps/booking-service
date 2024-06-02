package startup

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var unavailabilities = []*domain.Unavailability{
	{
		Id:              getObjectId("1643a56c9dea1760db469b7b"),
		AccommodationId: getObjectId("6643a56c9dea1760db469b7b"),
		UnavailabilityPeriods: []domain.UnavailabilityPeriod{
			{
				Id:     getObjectId("2643a56c9dea1760db469b7b"),
				Start:  time.Now(),
				End:    time.Now().Add(24 * time.Hour),
				Reason: domain.Reserved,
			},
			{
				Id:     getObjectId("3643a56c9dea1760db469b7b"),
				Start:  time.Now().Add(25 * time.Hour),
				End:    time.Now().Add(49 * time.Hour),
				Reason: domain.OwnerSet,
			},
		},
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
