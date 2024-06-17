package api

import (
	"context"
	"fmt"
	"github.com/ZMS-DevOps/booking-service/application"
	pb "github.com/ZMS-DevOps/booking-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookingHandler struct {
	pb.UnimplementedBookingServiceServer
	unavailabilityService     *application.UnavailabilityService
	reservationRequestService *application.ReservationRequestService
}

func NewBookingHandler(unavailabilityService *application.UnavailabilityService, reservationRequestService *application.ReservationRequestService) *BookingHandler {
	return &BookingHandler{
		unavailabilityService:     unavailabilityService,
		reservationRequestService: reservationRequestService,
	}
}

func (handler *BookingHandler) AddUnavailability(ctx context.Context, request *pb.AddUnavailabilityRequest) (*pb.AddUnavailabilityResponse, error) {
	id := request.Id
	accommodationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := handler.unavailabilityService.AddUnavailability(accommodationId, request.AccommodationName, request.Automatically, request.HostId); err != nil {
		return nil, err
	}
	return &pb.AddUnavailabilityResponse{}, nil
}

func (handler *BookingHandler) EditAccommodation(ctx context.Context, request *pb.EditAccommodationRequest) (*pb.EditAccommodationResponse, error) {
	id := request.Id
	accommodationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := handler.unavailabilityService.UpdateUnavailability(accommodationId, request.AccommodationName, request.Automatically, request.HostId); err != nil {
		return nil, err
	}
	return &pb.EditAccommodationResponse{}, nil
}

func (handler *BookingHandler) FilterAvailableAccommodation(ctx context.Context, request *pb.FilterAvailableAccommodationRequest) (*pb.FilterAvailableAccommodationResponse, error) {
	objectIDs, err := convertHexToObjectIDs(request.AccommodationIds)
	if err != nil {
		return nil, err
	}

	startDate, endDate, err := parseDates(request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}

	available, err := handler.unavailabilityService.FilterAvailable(objectIDs, startDate, endDate)
	if err != nil {
		return nil, err
	}

	accommodationIDs := make([]string, len(available))
	for i, id := range available {
		accommodationIDs[i] = id.Hex()
	}

	return &pb.FilterAvailableAccommodationResponse{AccommodationIds: accommodationIDs}, nil
}

func (handler *BookingHandler) CheckDeleteHost(ctx context.Context, request *pb.CheckDeleteHostRequest) (*pb.CheckDeleteHostResponse, error) {
	hostId := request.HostId
	success, err := handler.unavailabilityService.DeleteHost(hostId)
	if err != nil {
		return nil, err
	}
	return &pb.CheckDeleteHostResponse{Success: success}, nil
}

func (handler *BookingHandler) CheckDeleteClient(ctx context.Context, request *pb.CheckDeleteClientRequest) (*pb.CheckDeleteClientResponse, error) {
	clientId := request.HostId
	success := handler.reservationRequestService.DeleteClient(clientId)
	return &pb.CheckDeleteClientResponse{Success: success}, nil
}

func (handler *BookingHandler) CheckGuestHasReservationForHost(ctx context.Context, request *pb.CheckGuestHasReservationForHostRequest) (*pb.CheckGuestHasReservationForHostResponse, error) {
	reviewerId := request.ReviewerId
	hostId := request.HostId
	hasReservation := handler.reservationRequestService.CheckGuestHasReservationForHost(reviewerId, hostId)
	return &pb.CheckGuestHasReservationForHostResponse{HasReservation: hasReservation}, nil
}

func (handler *BookingHandler) CheckGuestHasReservationForAccommodation(ctx context.Context, request *pb.CheckGuestHasReservationForAccommodationRequest) (*pb.CheckGuestHasReservationForAccommodationResponse, error) {
	reviewerId := request.ReviewerId

	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		return nil, err
	}
	hasReservation := handler.reservationRequestService.CheckGuestHasReservationForAccommodation(reviewerId, accommodationId)
	return &pb.CheckGuestHasReservationForAccommodationResponse{HasReservation: hasReservation}, nil
}

func (handler *BookingHandler) CheckAccommodationHasReservation(ctx context.Context, request *pb.CheckAccommodationHasReservationRequest) (*pb.CheckAccommodationHasReservationResponse, error) {
	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		return nil, err
	}
	canDelete := handler.reservationRequestService.CheckAccommodationHasReservation(accommodationId)
	return &pb.CheckAccommodationHasReservationResponse{Success: canDelete}, nil
}

func convertHexToObjectIDs(hexIDs []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID

	for _, id := range hexIDs {
		idObj, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, idObj)
	}

	return objectIDs, nil
}

func parseDates(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	layout := time.RFC3339

	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing start date: %w", err)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error parsing end date: %w", err)
	}

	return startDate, endDate, nil
}
