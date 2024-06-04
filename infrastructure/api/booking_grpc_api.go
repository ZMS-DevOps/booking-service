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
	service *application.UnavailabilityService
}

func NewBookingHandler(service *application.UnavailabilityService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

func (handler *BookingHandler) AddUnavailability(ctx context.Context, request *pb.AddUnavailabilityRequest) (*pb.AddUnavailabilityResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := handler.service.AddUnavailability(objectId, request.AccommodationName, request.Automatically, request.HostId); err != nil {
		return nil, err
	}
	return &pb.AddUnavailabilityResponse{}, nil
}

func (handler *BookingHandler) EditAccommodation(ctx context.Context, request *pb.EditAccommodationRequest) (*pb.EditAccommodationResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := handler.service.UpdateUnavailability(objectId, request.AccommodationName, request.Automatically, request.HostId); err != nil {
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

	available, err := handler.service.FilterAvailable(objectIDs, startDate, endDate)
	if err != nil {
		return nil, err
	}

	accommodationIDs := make([]string, len(available))
	for i, id := range available {
		accommodationIDs[i] = id.Hex()
	}

	return &pb.FilterAvailableAccommodationResponse{AccommodationIds: accommodationIDs}, nil
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
