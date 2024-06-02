package api

import (
	"encoding/json"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	router.HandleFunc("/booking/request/{id}", handler.GetAll).Methods("GET")
	router.HandleFunc("/booking/request/{id}/pending", handler.GetPending).Methods("GET")
	router.HandleFunc("/booking/request/{id}/declined", handler.GetDeclined).Methods("GET")
	router.HandleFunc("/booking/request/{id}/approved", handler.GetApproved).Methods("GET")
	router.HandleFunc("/booking/request/{id}/completed", handler.GetCompleted).Methods("GET")
	//router.HandleFunc("/booking/request/{id}/approve", handler.Approve).Methods("PUT")
	//router.HandleFunc("/booking/request/{id}/decline", handler.Decline).Methods("PUT")
	//router.HandleFunc("/booking/request/{id}", handler.Delete).Methods("DELETE")
}

func (handler *ReservationRequestHandler) AddRequest(w http.ResponseWriter, r *http.Request) {
	var addReservationRequestDto dto.AddReservationRequestDto
	if err := json.NewDecoder(r.Body).Decode(&addReservationRequestDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateAddRegistrationRequestDto(addReservationRequestDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	newReservationRequest := dto.MapRegistrationRequest(addReservationRequestDto)

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

	response := dto.MapReservationRequestResponse(requests)

	jsonResponse, err := json.Marshal(response)
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

func (handler *ReservationRequestHandler) GetDeclined(w http.ResponseWriter, r *http.Request) {
	status := domain.Declined
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
