package unavailability

import (
	"context"
	"errors"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "bookingdb"
	COLLECTION = "unavailability"
)

type UnavailabilityMongoDBStore struct {
	unavailability *mongo.Collection
}

func NewUnavailabilityMongoDBStore(client *mongo.Client) domain.UnavailabilityStore {
	unavailability := client.Database(DATABASE).Collection(COLLECTION)
	return &UnavailabilityMongoDBStore{
		unavailability: unavailability,
	}
}

func (store *UnavailabilityMongoDBStore) Get(id primitive.ObjectID) (*domain.Unavailability, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UnavailabilityMongoDBStore) Insert(unavailability *domain.Unavailability) error {
	unavailability.Id = primitive.NewObjectID()
	result, err := store.unavailability.InsertOne(context.TODO(), unavailability)
	if err != nil {
		return err
	}
	unavailability.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *UnavailabilityMongoDBStore) DeleteAll() {
	store.unavailability.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *UnavailabilityMongoDBStore) GetAll() ([]*domain.Unavailability, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *UnavailabilityMongoDBStore) GetPeriod(periodId primitive.ObjectID) (domain.UnavailabilityPeriod, error) {
	var period domain.UnavailabilityPeriod

	filter := bson.M{"unavailability_periods._id": periodId}
	projection := bson.M{
		"unavailability_periods": bson.M{
			"$elemMatch": bson.M{"_id": periodId},
		},
	}
	var result domain.Unavailability
	err := store.unavailability.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return period, errors.New("period not found")
		}
		return period, err
	}
	if len(result.UnavailabilityPeriods) > 0 {
		period = result.UnavailabilityPeriods[0]
	} else {
		return period, errors.New("period not found")
	}

	return period, nil
}

func (store *UnavailabilityMongoDBStore) GetUnavailabilityPeriods(id primitive.ObjectID) ([]domain.UnavailabilityPeriod, error) {
	var unavailability domain.Unavailability
	filter := bson.M{"_id": id}
	err := store.unavailability.FindOne(context.TODO(), filter).Decode(&unavailability)
	if err != nil {
		return nil, err
	}
	return unavailability.UnavailabilityPeriods, nil
}

func (store *UnavailabilityMongoDBStore) UpdateUnavailabilityPeriods(unavailabilityId primitive.ObjectID, periods []domain.UnavailabilityPeriod) error {
	filter := bson.M{"_id": unavailabilityId}
	update := bson.M{}

	_, err := store.GetUnavailabilityPeriods(unavailabilityId)
	if err != nil {
		return err
	}

	update = bson.M{"$set": bson.M{"unavailability_periods": periods}}

	_, err = store.unavailability.UpdateOne(context.TODO(), filter, update)
	return err
}

func (store *UnavailabilityMongoDBStore) filterOne(filter interface{}) (unavailability *domain.Unavailability, err error) {
	result := store.unavailability.FindOne(context.TODO(), filter)
	err = result.Decode(&unavailability)
	return
}

func (store *UnavailabilityMongoDBStore) filter(filter interface{}) ([]*domain.Unavailability, error) {
	cursor, err := store.unavailability.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *UnavailabilityMongoDBStore) GetByAccommodationId(accommodationId primitive.ObjectID) (*domain.Unavailability, error) {
	var unavailability domain.Unavailability
	filter := bson.M{"accommodation_id": accommodationId}
	err := store.unavailability.FindOne(context.TODO(), filter).Decode(&unavailability)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &unavailability, nil
}

func (store *UnavailabilityMongoDBStore) GetByHostId(hostId string) ([]*domain.Unavailability, error) {
	filter := bson.M{"host_id": hostId}
	return store.filter(filter)
}

func (store *UnavailabilityMongoDBStore) Update(id primitive.ObjectID, unavailability *domain.Unavailability) error {
	filter := bson.M{"_id": id}

	updateFields := bson.D{
		{"accommodation_id", unavailability.AccommodationId},
		{"unavailability_periods", unavailability.UnavailabilityPeriods},
		{"review_reservation_request_automatically", unavailability.ReviewReservationRequestAutomatically},
	}
	update := bson.D{{"$set", updateFields}}

	_, err := store.unavailability.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func decode(cursor *mongo.Cursor) (unavailabilities []*domain.Unavailability, err error) {
	for cursor.Next(context.TODO()) {
		var unavailability domain.Unavailability
		err = cursor.Decode(&unavailability)
		if err != nil {
			return
		}
		unavailabilities = append(unavailabilities, &unavailability)
	}
	err = cursor.Err()
	return
}
