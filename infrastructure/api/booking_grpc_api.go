package api

import (
	"context"
	"fmt"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	pb "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/booking-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

type BookingHandler struct {
	pb.UnimplementedBookingServiceServer
	unavailabilityService     *application.UnavailabilityService
	reservationRequestService *application.ReservationRequestService
	traceProvider             *sdktrace.TracerProvider
	loki                      promtail.Client
}

func NewBookingHandler(unavailabilityService *application.UnavailabilityService, reservationRequestService *application.ReservationRequestService, provider *sdktrace.TracerProvider, loki promtail.Client) *BookingHandler {
	return &BookingHandler{
		unavailabilityService:     unavailabilityService,
		reservationRequestService: reservationRequestService,
		traceProvider:             provider,
		loki:                      loki,
	}
}

func (handler *BookingHandler) AddUnavailability(ctx context.Context, request *pb.AddUnavailabilityRequest) (*pb.AddUnavailabilityResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "add-unavailability-grpc")
	defer func() { span.End() }()
	id := request.Id
	accommodationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "AddUnavailability", "")
		return nil, err
	}
	if err := handler.unavailabilityService.AddUnavailability(accommodationId, request.AccommodationName, request.Automatically, request.HostId, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to add unavailability", span, handler.loki, "AddUnavailability", "")
		return nil, err
	}
	util.HttpTraceInfo("Unavailability added successfully", span, handler.loki, "AddUnavailability", "")
	return &pb.AddUnavailabilityResponse{}, nil
}

func (handler *BookingHandler) EditAccommodation(ctx context.Context, request *pb.EditAccommodationRequest) (*pb.EditAccommodationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "edit-accommodation-grpc")
	defer func() { span.End() }()
	id := request.Id
	accommodationId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "EditAccommodation", "")
		return nil, err
	}
	if err := handler.unavailabilityService.UpdateUnavailability(accommodationId, request.AccommodationName, request.Automatically, request.HostId, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to update unavailability", span, handler.loki, "EditAccommodation", "")
		return nil, err
	}
	util.HttpTraceInfo("Edit accommodation successful", span, handler.loki, "EditAccommodation", "")
	return &pb.EditAccommodationResponse{}, nil
}

func (handler *BookingHandler) FilterAvailableAccommodation(ctx context.Context, request *pb.FilterAvailableAccommodationRequest) (*pb.FilterAvailableAccommodationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "filter-available-accommodation-grpc")
	defer func() { span.End() }()
	objectIDs, err := convertHexToObjectIDs(request.AccommodationIds)
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "FilterAvailableAccommodation", "")
		return nil, err
	}

	startDate, endDate, err := parseDates(request.StartDate, request.EndDate)
	if err != nil {
		util.HttpTraceError(err, "failed to parse dates", span, handler.loki, "FilterAvailableAccommodation", "")
		return nil, err
	}

	available, err := handler.unavailabilityService.FilterAvailable(objectIDs, startDate, endDate, span)
	if err != nil {
		util.HttpTraceError(err, "failed to filter available accommodation", span, handler.loki, "FilterAvailableAccommodation", "")
		return nil, err
	}

	accommodationIDs := make([]string, len(available))
	for i, id := range available {
		accommodationIDs[i] = id.Hex()
	}

	util.HttpTraceInfo("Accommodations filtered successfully", span, handler.loki, "FilterAvailableAccommodation", "")
	return &pb.FilterAvailableAccommodationResponse{AccommodationIds: accommodationIDs}, nil
}

func (handler *BookingHandler) CheckDeleteHost(ctx context.Context, request *pb.CheckDeleteHostRequest) (*pb.CheckDeleteHostResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "check-delete-host-grpc")
	defer func() { span.End() }()
	hostId := request.HostId
	success, err := handler.unavailabilityService.DeleteHost(hostId, span, handler.loki)
	if err != nil {
		util.HttpTraceError(err, "failed to delete unavailability", span, handler.loki, "CheckDeleteHost", "")
		return nil, err
	}
	util.HttpTraceInfo("Check delete host processed successfully", span, handler.loki, "CheckDeleteHost", "")
	return &pb.CheckDeleteHostResponse{Success: success}, nil
}

func (handler *BookingHandler) CheckDeleteClient(ctx context.Context, request *pb.CheckDeleteClientRequest) (*pb.CheckDeleteClientResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "check-delete-client-grpc")
	defer func() { span.End() }()
	clientId := request.HostId
	success := handler.reservationRequestService.DeleteClient(clientId, span, handler.loki)
	util.HttpTraceInfo("Check delete client processed successfully", span, handler.loki, "CheckDeleteClient", "")
	return &pb.CheckDeleteClientResponse{Success: success}, nil
}

func (handler *BookingHandler) CheckGuestHasReservationForHost(ctx context.Context, request *pb.CheckGuestHasReservationForHostRequest) (*pb.CheckGuestHasReservationForHostResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "check-guest-has-reservation-request-for-host-grpc")
	defer func() { span.End() }()
	reviewerId := request.ReviewerId
	hostId := request.HostId
	hasReservation := handler.reservationRequestService.CheckGuestHasReservationForHost(reviewerId, hostId, span, handler.loki)
	util.HttpTraceInfo("Check quest has reservation for host processed successfully", span, handler.loki, "CheckGuestHasReservationForHost", "")
	return &pb.CheckGuestHasReservationForHostResponse{HasReservation: hasReservation}, nil
}

func (handler *BookingHandler) CheckGuestHasReservationForAccommodation(ctx context.Context, request *pb.CheckGuestHasReservationForAccommodationRequest) (*pb.CheckGuestHasReservationForAccommodationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "check-guest-has-reservation-for-accommodation-grpc")
	defer func() { span.End() }()
	reviewerId := request.ReviewerId

	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "CheckGuestHasReservationForAccommodation", "")
		return nil, err
	}
	hasReservation := handler.reservationRequestService.CheckGuestHasReservationForAccommodation(reviewerId, accommodationId, span, handler.loki)
	util.HttpTraceInfo("Check quest has reservation for accommodation processed successfully", span, handler.loki, "CheckGuestHasReservationForAccommodation", "")
	return &pb.CheckGuestHasReservationForAccommodationResponse{HasReservation: hasReservation}, nil
}

func (handler *BookingHandler) CheckAccommodationHasReservation(ctx context.Context, request *pb.CheckAccommodationHasReservationRequest) (*pb.CheckAccommodationHasReservationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "check-accommodation-has-reservation-grpc")
	defer func() { span.End() }()
	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "CheckAccommodationHasReservation", "")
		return nil, err
	}
	canDelete := handler.reservationRequestService.CheckAccommodationHasReservation(accommodationId, span, handler.loki)
	util.HttpTraceInfo("Check quest has reservation processed successfully", span, handler.loki, "CheckAccommodationHasReservation", "")
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
