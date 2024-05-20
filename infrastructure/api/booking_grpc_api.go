package api

import (
	"context"
	"fmt"
	"github.com/ZMS-DevOps/booking-service/application"
	pb "github.com/ZMS-DevOps/booking-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	fmt.Println("Stigao i jaaa")
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	if err := handler.service.AddUnavailability(objectId); err != nil {
		return nil, err
	}
	return &pb.AddUnavailabilityResponse{}, nil
}
