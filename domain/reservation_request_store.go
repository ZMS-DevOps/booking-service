package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReservationRequestStore interface {
	Get(id primitive.ObjectID) (*ReservationRequest, error)
	Insert(reservationRequest *ReservationRequest) error
	Update(id primitive.ObjectID, reservationRequest *ReservationRequest) error
	DeleteAll()
	GetAll() ([]*ReservationRequest, error)
	GetByAccommodationId(accommodationId primitive.ObjectID) ([]*ReservationRequest, error)
	GetByAccommodationIdAndType(id primitive.ObjectID, requestType ReservationRequestStatus) ([]*ReservationRequest, error)
	Delete(id primitive.ObjectID) error
	CancelOverlappingPendingRequests(request *ReservationRequest) error
	GetByClientId(clientId primitive.ObjectID) ([]*ReservationRequest, error)
}
