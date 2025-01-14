package reservation_request

import (
	"context"
	"github.com/ZMS-DevOps/booking-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
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

func (store *ReservationRequestMongoDBStore) DeleteAll() {
	store.reservationRequestCollection.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *ReservationRequestMongoDBStore) Insert(reservationRequest *domain.ReservationRequest) (*primitive.ObjectID, error) {
	reservationRequest.Id = primitive.NewObjectID()
	result, err := store.reservationRequestCollection.InsertOne(context.TODO(), reservationRequest)
	if err != nil {
		return nil, err
	}
	reservationRequest.Id = result.InsertedID.(primitive.ObjectID)
	return &reservationRequest.Id, nil
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

func (store *ReservationRequestMongoDBStore) Update(id primitive.ObjectID, reservationRequest *domain.ReservationRequest) error {
	filter := bson.M{"_id": id}

	updateFields := bson.D{
		{"accommodation_id", reservationRequest.AccommodationId},
		{"accommodation_name", reservationRequest.AccommodationName},
		{"user_id", reservationRequest.UserId},
		{"start", reservationRequest.Start},
		{"end", reservationRequest.End},
		{"number_of_guests", reservationRequest.NumberOfGuests},
		{"price_total", reservationRequest.PriceTotal},
		{"status", reservationRequest.Status},
	}
	update := bson.D{{"$set", updateFields}}

	_, err := store.reservationRequestCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (store *ReservationRequestMongoDBStore) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := store.reservationRequestCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *ReservationRequestMongoDBStore) DeleteByHost(id string) error {
	filter := bson.M{"host_id": id}
	_, err := store.reservationRequestCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *ReservationRequestMongoDBStore) CancelOverlappingPendingRequests(reservationRequest *domain.ReservationRequest) error {
	filter := bson.M{
		"accommodation_id": reservationRequest.AccommodationId,
		"status":           0,
		"start":            bson.M{"$lt": reservationRequest.End},
		"end":              bson.M{"$gt": reservationRequest.Start},
	}

	update := bson.M{
		"$set": bson.M{
			"status": domain.DeclinedByHost,
		},
	}

	updateResult, err := store.reservationRequestCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if updateResult.MatchedCount > 0 {
		log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	} else {
		log.Printf("No documents matched the query.")
	}

	return nil
}

func (store *ReservationRequestMongoDBStore) GetByClientId(clientId string) ([]*domain.ReservationRequest, error) {
	filter := bson.M{"user_id": clientId}
	return store.filter(filter)
}

func (store *ReservationRequestMongoDBStore) GetByClientIdAndStatus(userId string, status domain.ReservationRequestStatus) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"user_id": userId,
		"status":  status,
	}
	return store.filter(filter)
}

func (store *ReservationRequestMongoDBStore) GetByHostAndTimeAndSearch(hostId string, past bool, search string) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"host_id": hostId,
	}
	now := time.Now()
	if past {
		filter["end"] = bson.M{"$lt": now}
	} else {
		filter["end"] = bson.M{"$gte": now}
	}

	store.searchByString(search, filter)

	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return decodeReservationRequests(cursor)
}

func (store *ReservationRequestMongoDBStore) searchByString(search string, filter bson.M) {
	if search != "" {
		filter["$or"] = []bson.M{
			{"accommodation_name": bson.M{"$regex": search, "$options": "i"}},
			{"status": bson.M{"$regex": search, "$options": "i"}},
		}
	}
}

func (store *ReservationRequestMongoDBStore) GetByClientIdAndTimeAndSearch(guestId string, past bool, search string) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"user_id": guestId,
	}
	now := time.Now()
	if past {
		filter["end"] = bson.M{"$lt": now}
	} else {
		filter["end"] = bson.M{"$gte": now}
	}

	store.searchByString(search, filter)

	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return decodeReservationRequests(cursor)
}

func (store *ReservationRequestMongoDBStore) GetByClientIdAndHostId(clientId string, hostId string) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"user_id": clientId,
		"host_id": hostId,
	}

	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return decodeReservationRequests(cursor)
}

func (store *ReservationRequestMongoDBStore) GetPastAcceptedReservationRequestByClientIdAndHostId(clientId string, hostId string) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"user_id": clientId,
		"host_id": hostId,
		"status":  1,
		"end": bson.M{
			"$lt": time.Now(),
		},
	}

	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return decodeReservationRequests(cursor)
}

func (store *ReservationRequestMongoDBStore) GetByClientIdAndAccommodationId(reviewerId string, accommodationId primitive.ObjectID) ([]*domain.ReservationRequest, error) {
	filter := bson.M{
		"user_id":          reviewerId,
		"accommodation_id": accommodationId,
		"status":           1,
		"end": bson.M{
			"$lt": time.Now(),
		},
	}

	cursor, err := store.reservationRequestCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return decodeReservationRequests(cursor)
}

func (store *ReservationRequestMongoDBStore) DeleteByAccommodation(accommodationId primitive.ObjectID) error {
	filter := bson.M{
		"status":           0,
		"accommodation_id": accommodationId,
	}
	update := bson.M{
		"$set": bson.M{
			"status": domain.DeclinedByHost,
		},
	}

	_, err := store.reservationRequestCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
