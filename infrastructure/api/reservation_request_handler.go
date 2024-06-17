package api

import (
	"encoding/json"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type ReservationRequestHandler struct {
	service *application.ReservationRequestService
}

func NewReservationRequestHandler(service *application.ReservationRequestService) *ReservationRequestHandler {
	server := &ReservationRequestHandler{
		service: service,
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
	var addReservationRequestDto dto.AddReservationRequestDto
	if err := json.NewDecoder(r.Body).Decode(&addReservationRequestDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Printf("prvi print")
	log.Printf("adddto print %v", addReservationRequestDto.AccommodationId)

	if err := dto.ValidateAddRegistrationRequestDto(addReservationRequestDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("drugi print")

	newReservationRequest := dto.MapRegistrationRequest(addReservationRequestDto)
	log.Printf("newReservationRequest print %v", newReservationRequest)

	if err := handler.service.AddReservationRequest(newReservationRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ReservationRequestHandler) GetByStatus(w http.ResponseWriter, r *http.Request, status *domain.ReservationRequestStatus) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requests, err := handler.service.GetByAccommodationId(accommodationId, status)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responses := dto.MapReservationRequestResponse(requests)

	for _, response := range responses {
		response.NumberOfCanceledReservations = handler.service.GetNumberOfCanceled(response.UserId)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *ReservationRequestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	handler.GetByStatus(w, r, nil)
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
	vars := mux.Vars(r)
	reservationRequestId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.ApproveRequest(reservationRequestId); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ReservationRequestHandler) Decline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reservationRequestId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.DeclineRequest(reservationRequestId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ReservationRequestHandler) DeclineReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.service.DeclineReservation(accommodationId); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ReservationRequestHandler) GetByClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	requests, err := handler.service.GetByClientId(userId, nil)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responses := dto.MapReservationRequestResponse(requests)

	for _, response := range responses {
		response.NumberOfCanceledReservations = handler.service.GetNumberOfCanceled(response.UserId)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *ReservationRequestHandler) GetFilteredRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	userType := r.URL.Query().Get("user-type")
	if userType != "host" && userType != "guest" {
		handleError(w, http.StatusBadRequest, "Invalid user type")
		return
	}

	pastStr := r.URL.Query().Get("past")
	past := false
	if pastStr == "true" {
		past = true
	} else if pastStr != "false" {
		handleError(w, http.StatusBadRequest, "Invalid past parameter")
		return
	}

	searchStr := r.URL.Query().Get("search")

	requests, err := handler.service.GetFilteredRequests(userId, userType, past, searchStr)
	if err != nil {
		handleError(w, http.StatusNotFound, "No reservations found")
		return
	}

	responses := dto.MapReservationRequestResponse(requests)
	for _, response := range responses {
		response.NumberOfCanceledReservations = handler.service.GetNumberOfCanceled(response.UserId)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Error marshalling response")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
