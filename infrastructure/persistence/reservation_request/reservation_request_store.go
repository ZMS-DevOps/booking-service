package reservation_request

import (
	"context"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	DATABASE   = "bookingdb"
	COLLECTION = "reservation_request"
)

type ReservationRequestMongoDBStore struct {
	reservationRequestCollection *mongo.Collection
}

func NewReservationRequestMongoDBStore(client *mongo.Client) domain.ReservationRequestStore {
	reservationRequestCollection := client.Database(DATABASE).Collection(COLLECTION)
	return &ReservationRequestMongoDBStore{
		reservationRequestCollection: reservationRequestCollection,
	}
}

func (store *ReservationRequestMongoDBStore) Get(id primitive.ObjectID) (*domain.ReservationRequest, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *ReservationRequestMongoDBStore) Insert(reservationRequest *domain.ReservationRequest) error {
	reservationRequest.Id = primitive.NewObjectID()
	result, err := store.reservationRequestCollection.InsertOne(context.TODO(), reservationRequest)
	if err != nil {
		return err
	}
	reservationRequest.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *ReservationRequestMongoDBStore) DeleteAll() {
	store.reservationRequestCollection.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *ReservationRequestMongoDBStore) GetAll() ([]*domain.ReservationRequest, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *ReservationRequestMongoDBStore) filterOne(filter interface{}) (reservationRequest *domain.ReservationRequest, err error) {
	result := store.reservationRequestCollection.FindOne(context.TODO(), filter)
	err = result.Decode(&reservationRequest)
	return
}

func (store *ReservationRequestMongoDBStore) filter(filter interface{}) ([]*domain.ReservationRequest, error) {
	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decodeReservationRequests(cursor)
}

func decodeReservationRequests(cursor *mongo.Cursor) (reservationRequests []*domain.ReservationRequest, err error) {
	for cursor.Next(context.TODO()) {
		var reservationRequest domain.ReservationRequest
		err = cursor.Decode(&reservationRequest)
		if err != nil {
			return
		}
		reservationRequests = append(reservationRequests, &reservationRequest)
	}
	err = cursor.Err()
	return
}

func (store *ReservationRequestMongoDBStore) GetByAccommodationId(accommodationId primitive.ObjectID) ([]*domain.ReservationRequest, error) {
	filter := bson.M{"accommodation_id": accommodationId}
	var reservationRequests []*domain.ReservationRequest
	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &reservationRequests); err != nil {
		log.Fatal(err)
	}

	if err := cursor.Close(context.TODO()); err != nil {
		log.Fatal(err)
	}
	return reservationRequests, nil
}

func (store *ReservationRequestMongoDBStore) GetByAccommodationIdAndType(accommodationId primitive.ObjectID, requestType domain.ReservationRequestStatus) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"accommodation_id": accommodationId,
		"status":           requestType,
	}
	var reservationRequests []*domain.ReservationRequest
	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &reservationRequests); err != nil {
		log.Fatal(err)
	}

	if err := cursor.Close(context.TODO()); err != nil {
		log.Fatal(err)
	}
	return reservationRequests, nil
}