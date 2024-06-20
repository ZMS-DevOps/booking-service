package api

import (
	"encoding/json"
	"errors"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/ZMS-DevOps/booking-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
)

type ReservationRequestHandler struct {
	service       *application.ReservationRequestService
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

func NewReservationRequestHandler(service *application.ReservationRequestService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *ReservationRequestHandler {
	server := &ReservationRequestHandler{
		service:       service,
		traceProvider: traceProvider,
		loki:          loki,
	}
	return server
}

func (handler *ReservationRequestHandler) Init(router *mux.Router) {
	router.HandleFunc("/booking/request", handler.AddRequest).Methods("POST")
	router.HandleFunc("/booking/reservation/decline/{id}", handler.DeclineReservation).Methods("PUT")
	router.HandleFunc("/booking/request/user/{id}", handler.GetFilteredRequests).Methods("GET")
	router.HandleFunc("/booking/request/all/{id}", handler.GetAll).Methods("GET")
	router.HandleFunc("/booking/request/approve/{id}", handler.Approve).Methods("PUT")
	router.HandleFunc("/booking/request/decline/{id}", handler.Decline).Methods("PUT")
}

func (handler *ReservationRequestHandler) AddRequest(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "add-request-post")
	defer func() { span.End() }()
	var addReservationRequestDto dto.AddReservationRequestDto
	if err := json.NewDecoder(r.Body).Decode(&addReservationRequestDto); err != nil {
		util.HttpTraceError(err, "invalid request payload", span, handler.loki, "AddRequest", "")
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Printf("prvi print")
	log.Printf("adddto print %v", addReservationRequestDto.AccommodationId)

	if err := dto.ValidateAddRegistrationRequestDto(addReservationRequestDto); err != nil {
		util.HttpTraceError(err, "failed to validate request body", span, handler.loki, "AddRequest", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("drugi print")

	newReservationRequest := dto.MapRegistrationRequest(addReservationRequestDto)
	log.Printf("newReservationRequest print %v", newReservationRequest)

	if err := handler.service.AddReservationRequest(newReservationRequest, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to add reservation request", span, handler.loki, "AddRequest", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	util.HttpTraceInfo("Reservation request added successfully", span, handler.loki, "AddRequest", "")
	w.WriteHeader(http.StatusCreated)
}

func (handler *ReservationRequestHandler) GetByStatus(w http.ResponseWriter, r *http.Request, status *domain.ReservationRequestStatus) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-status")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "failed to get accommodation by id", span, handler.loki, "GetByStatus", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requests, err := handler.service.GetByAccommodationId(accommodationId, status, span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "failed to get accommodation by id and status", span, handler.loki, "GetByStatus", "")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responses := dto.MapReservationRequestResponse(requests)

	for _, response := range responses {
		response.NumberOfCanceledReservations = handler.service.GetNumberOfCanceled(response.UserId, span, handler.loki)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetByStatus", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Reservation request fetched by status successfully", span, handler.loki, "GetByStatus", "")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *ReservationRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-all-get")
	defer func() { span.End() }()
	handler.GetByStatus(w, r, nil)
	util.HttpTraceInfo("Reservation requests fetched successfully", span, handler.loki, "GetAll", "")
}

func (handler *ReservationRequestHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	status := domain.Pending
	handler.GetByStatus(w, r, &status)
}

func (handler *ReservationRequestHandler) GetApproved(w http.ResponseWriter, r *http.Request) {
	status := domain.Approved
	handler.GetByStatus(w, r, &status)
}

func (handler *ReservationRequestHandler) GetCompleted(w http.ResponseWriter, r *http.Request) {
	status := domain.Completed
	handler.GetByStatus(w, r, &status)
}

func (handler *ReservationRequestHandler) Approve(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "approve-put")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	reservationRequestId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "failed to get reservation request by id", span, handler.loki, "Approve", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.ApproveRequest(reservationRequestId, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to approve request", span, handler.loki, "Approve", "")
		w.WriteHeader(http.StatusNotFound)
	}

	util.HttpTraceInfo("Reservation request approved successfully", span, handler.loki, "Approve", "")
	w.WriteHeader(http.StatusOK)
}

func (handler *ReservationRequestHandler) Decline(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "add-post")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	reservationRequestId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "invalid reservation request id", span, handler.loki, "Decline", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.DeclineRequest(reservationRequestId, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to decline request", span, handler.loki, "Decline", "")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	util.HttpTraceInfo("Reservation request successfully declined", span, handler.loki, "Decline", "")
	w.WriteHeader(http.StatusOK)
}

func (handler *ReservationRequestHandler) DeclineReservation(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "decline-reservation-put")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "DeclineReservation", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.DeclineReservation(accommodationId, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to decline reservation", span, handler.loki, "DeclineReservation", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	util.HttpTraceInfo("Reservation was successfully declined", span, handler.loki, "DeclineReservation", "")
	w.WriteHeader(http.StatusOK)
}

func (handler *ReservationRequestHandler) GetByClient(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-client-get")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	userId := vars["id"]
	requests, err := handler.service.GetByClientId(userId, nil, span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "invalid reservation request", span, handler.loki, "GetByClient", "")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responses := dto.MapReservationRequestResponse(requests)

	for _, response := range responses {
		response.NumberOfCanceledReservations = handler.service.GetNumberOfCanceled(response.UserId, span, handler.loki)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		util.HttpTraceError(err, "failed marshal data", span, handler.loki, "GetByClient", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Reservation request fetched successfully", span, handler.loki, "GetByClient", "")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *ReservationRequestHandler) GetFilteredRequests(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-filtered-requests-get")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	userId := vars["id"]
	userType := r.URL.Query().Get("user-type")
	if userType != "host" && userType != "guest" {
		util.HttpTraceError(errors.New("invalid user type"), "invalid user type", span, handler.loki, "GetFilteredRequests", "")
		handleError(w, http.StatusBadRequest, "Invalid user type")
		return
	}

	pastStr := r.URL.Query().Get("past")
	past := false
	if pastStr == "true" {
		past = true
	} else if pastStr != "false" {
		util.HttpTraceError(errors.New("invalid past parameter"), "invalid past parameter", span, handler.loki, "GetFilteredRequests", "")
		handleError(w, http.StatusBadRequest, "Invalid past parameter")
		return
	}

	searchStr := r.URL.Query().Get("search")

	requests, err := handler.service.GetFilteredRequests(userId, userType, past, searchStr, span, handler.loki)
	if err != nil {
		util.HttpTraceError(err, "no reservations found", span, handler.loki, "GetFilteredRequests", "")
		handleError(w, http.StatusNotFound, "No reservations found")
		return
	}

	responses := dto.MapReservationRequestResponse(requests)
	for _, response := range responses {
		response.NumberOfCanceledReservations = handler.service.GetNumberOfCanceled(response.UserId, span, handler.loki)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetFilteredRequests", "")
		handleError(w, http.StatusInternalServerError, "Error marshalling response")
		return
	}
	util.HttpTraceInfo("Reservation request successfully filtered", span, handler.loki, "GetFilteredRequests", "")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
