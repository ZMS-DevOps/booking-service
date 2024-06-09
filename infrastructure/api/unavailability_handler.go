package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UnavailabilityHandler struct {
	service *application.UnavailabilityService
}

type HealthCheckResponse struct {
	Size string `json:"size"`
}

func NewUnavailabilityHandler(service *application.UnavailabilityService) *UnavailabilityHandler {
	server := &UnavailabilityHandler{
		service: service,
	}
	return server
}

func (handler *UnavailabilityHandler) Init(router *mux.Router) {
	router.HandleFunc("/booking/health", handler.GetHealthCheck).Methods("GET")
	router.HandleFunc("/booking/unavailability", handler.GetAll).Methods("GET")
	router.HandleFunc("/booking/unavailability/{id}", handler.GetByAccommodationId).Methods("GET")
	router.HandleFunc("/booking/unavailability/host/{id}", handler.GetByHostId).Methods("GET")
	router.HandleFunc("/booking/unavailability/remove", handler.DeletePeriod).Methods("PUT")
	router.HandleFunc("/booking/unavailability/add", handler.AddPeriod).Methods("PUT")
}

func (handler *UnavailabilityHandler) AddPeriod(w http.ResponseWriter, r *http.Request) {
	var manageUnavailabilityPeriodDto dto.ManageUnavailabilityPeriodDto
	if err := json.NewDecoder(r.Body).Decode(&manageUnavailabilityPeriodDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateManageUnavailabilityPeriodDto(manageUnavailabilityPeriodDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUnavailabilityPeriod := dto.MapUnavailabilityPeriod(&manageUnavailabilityPeriodDto)
	newUnavailabilityPeriod.Reason = domain.OwnerSet
	if err := handler.service.AddUnavailabilityPeriod(manageUnavailabilityPeriodDto.AccommodationId, newUnavailabilityPeriod); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *UnavailabilityHandler) DeletePeriod(w http.ResponseWriter, r *http.Request) {
	var manageUnavailabilityPeriodDto dto.ManageUnavailabilityPeriodDto
	if err := json.NewDecoder(r.Body).Decode(&manageUnavailabilityPeriodDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateManageUnavailabilityPeriodDto(manageUnavailabilityPeriodDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	removedUnavailabilityPeriod := dto.MapUnavailabilityPeriod(&manageUnavailabilityPeriodDto)
	if err := handler.service.RemoveUnavailabilityPeriod(manageUnavailabilityPeriodDto.AccommodationId, removedUnavailabilityPeriod); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *UnavailabilityHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	unavailabilityList, err := handler.service.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responseList []dto.UnavailabilityResponse
	for _, unavailability := range unavailabilityList {
		for _, period := range unavailability.UnavailabilityPeriods {
			response := dto.MapToUnavailabilityResponse(
				period.Id,
				unavailability.AccommodationId,
				unavailability.AccommodationName,
				period.Start,
				period.End,
				period.Reason,
			)
			responseList = append(responseList, response)
		}
	}

	jsonResponse, err := json.Marshal(responseList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *UnavailabilityHandler) GetByAccommodationId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unavailability, err := handler.service.GetByAccommodationId(accommodationId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responseList []dto.UnavailabilityResponse
	for _, period := range unavailability.UnavailabilityPeriods {
		response := dto.MapToUnavailabilityResponse(
			period.Id,
			unavailability.AccommodationId,
			unavailability.AccommodationName,
			period.Start,
			period.End,
			period.Reason,
		)
		responseList = append(responseList, response)
	}

	jsonResponse, err := json.Marshal(responseList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *UnavailabilityHandler) GetByHostId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unavailabilityList, err := handler.service.GetByHostId(accommodationId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var responseList []dto.UnavailabilityResponse
	for _, unavailability := range unavailabilityList {
		for _, period := range unavailability.UnavailabilityPeriods {
			response := dto.MapToUnavailabilityResponse(
				period.Id,
				unavailability.AccommodationId,
				unavailability.AccommodationName,
				period.Start,
				period.End,
				period.Reason,
			)
			responseList = append(responseList, response)
		}
	}

	jsonResponse, err := json.Marshal(responseList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *UnavailabilityHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Size: "Booking SERVICE OK",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}
