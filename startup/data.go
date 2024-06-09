package startup

import (
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var unavailabilities = []*domain.Unavailability{
	{
		Id:                getObjectId("1643a56c9dea1760db469b7b"),
		AccommodationId:   getObjectId("gg43bdc7240f80f13b5d24e9"),
		AccommodationName: "Balinese Bungalow",
		HostId:            getObjectId("57325353-5469-4930-8ec9-35c003e1b967"),
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
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                getObjectId("2743a56c9dea1760db469b7b"),
		AccommodationId:   getObjectId("hh43bdc7240f80f13b5d25e0"),
		AccommodationName: "Cliffside Retreat",
		HostId:            getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods: []domain.UnavailabilityPeriod{
			{
				Id:     getObjectId("3743a56c9dea1760db469b7b"),
				Start:  time.Now(),
				End:    time.Now().Add(24 * time.Hour),
				Reason: domain.Reserved,
			},
			{
				Id:     getObjectId("4743a56c9dea1760db469b7b"),
				Start:  time.Now().Add(25 * time.Hour),
				End:    time.Now().Add(49 * time.Hour),
				Reason: domain.OwnerSet,
			},
		},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("ff43bdc7240f80f13b5d23e8"),
		AccommodationName:                     "Luxury Villa",
		HostId:                                getObjectId("57325353-5469-4930-8ec9-35c003e1b967"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("gg43bdc7240f80f13b5d24e9"),
		AccommodationName:                     "Balinese Bungalow",
		HostId:                                getObjectId("57325353-5469-4930-8ec9-35c003e1b967"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("hh43bdc7240f80f13b5d25e0"),
		AccommodationName:                     "Cliffside Retreat",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("ii43bdc7240f80f13b5d26e1"),
		AccommodationName:                     "Ubud Villa",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("7743a6dc9dea1760db469b8c"),
		AccommodationName:                     "Mountain Retreat",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("8843bdc7240f80f13b5d19e7"),
		AccommodationName:                     "Urban Loft",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("jj43bdc7240f80f13b5d27e2"),
		AccommodationName:                     "Central Park Apartment",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("kk43bdc7240f80f13b5d28e3"),
		AccommodationName:                     "Brooklyn Brownstone",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("ll43bdc7240f80f13b5d29e4"),
		AccommodationName:                     "Times Square Studio",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("9943a7dc9dea1760db469b9d"),
		AccommodationName:                     "Beachfront Villa",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("aa43bdc7240f80f13b5d20d7"),
		AccommodationName:                     "Country Farmhouse",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("bb43a8dc9dea1760db469c1e"),
		AccommodationName:                     "Desert Oasis",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("cc43bdc7240f80f13b5d21e7"),
		AccommodationName:                     "Tropical Paradise",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("dd43a9dc9dea1760db469c2f"),
		AccommodationName:                     "Safari Lodge",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: false,
	},
	{
		Id:                                    primitive.NewObjectID(),
		AccommodationId:                       getObjectId("ee43bdc7240f80f13b5d22e7"),
		AccommodationName:                     "Lake House",
		HostId:                                getObjectId("04d19820-6340-4c93-84f9-2ffda959a0d9"),
		UnavailabilityPeriods:                 []domain.UnavailabilityPeriod{},
		ReviewReservationRequestAutomatically: true,
	},
	{
		Id:                getObjectId("2743a56c9dea1760db469b7b"),
		AccommodationId:   getObjectId("7743a56c9dea1760db469b7b"),
		AccommodationName: "Other accommodation name",
		HostId:            getObjectId("2243a56c9dea1760db469b7b"),
		UnavailabilityPeriods: []domain.UnavailabilityPeriod{
			{
				Id:     getObjectId("3743a56c9dea1760db469b7b"),
				Start:  time.Now(),
				End:    time.Now().Add(24 * time.Hour),
				Reason: domain.Reserved,
			},
			{
				Id:     getObjectId("4743a56c9dea1760db469b7b"),
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
