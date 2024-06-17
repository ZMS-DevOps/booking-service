package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReservationRequestStore interface {
	Get(id primitive.ObjectID) (*ReservationRequest, error)
	Insert(reservationRequest *ReservationRequest) (*primitive.ObjectID, error)
	Update(id primitive.ObjectID, reservationRequest *ReservationRequest) error
	DeleteAll()
	GetAll() ([]*ReservationRequest, error)
	GetByAccommodationId(accommodationId primitive.ObjectID) ([]*ReservationRequest, error)
	GetByAccommodationIdAndType(id primitive.ObjectID, requestType ReservationRequestStatus) ([]*ReservationRequest, error)
	Delete(id primitive.ObjectID) error
	CancelOverlappingPendingRequests(request *ReservationRequest) error
	GetByClientId(clientId string) ([]*ReservationRequest, error)
	GetByClientIdAndStatus(clientId string, status ReservationRequestStatus) ([]*ReservationRequest, error)
	GetByHostAndTimeAndSearch(userId string, past bool, search string) ([]*ReservationRequest, error)
	GetByClientIdAndTimeAndSearch(guestId string, past bool, search string) ([]*ReservationRequest, error)
	DeleteByHost(hostId string) error
	GetByClientIdAndHostId(reviewerId string, hostId string) ([]*ReservationRequest, error)
	GetPastAcceptedReservationRequestByClientIdAndHostId(reviewerId string, hostId string) ([]*ReservationRequest, error)
	GetByClientIdAndAccommodationId(reviewerId string, accommodationId primitive.ObjectID) ([]*ReservationRequest, error)
	DeleteByAccommodation(accommodationId primitive.ObjectID) error
}
