package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/booking-service/application"
	"github.com/ZMS-DevOps/booking-service/domain"
	"github.com/ZMS-DevOps/booking-service/infrastructure/dto"
	"github.com/ZMS-DevOps/booking-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

type UnavailabilityHandler struct {
	service       *application.UnavailabilityService
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

type HealthCheckResponse struct {
	Size string `json:"size"`
}

func NewUnavailabilityHandler(service *application.UnavailabilityService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *UnavailabilityHandler {
	server := &UnavailabilityHandler{
		service:       service,
		traceProvider: traceProvider,
		loki:          loki,
	}
	return server
}

func (handler *UnavailabilityHandler) Init(router *mux.Router) {
	router.HandleFunc("/booking/unavailability", handler.GetAll).Methods("GET")
	router.HandleFunc("/booking/unavailability/accommodation/{id}", handler.GetByAccommodationId).Methods("GET")
	router.HandleFunc("/booking/unavailability/host/{id}", handler.GetByHostId).Methods("GET")
	router.HandleFunc("/booking/unavailability/remove", handler.DeletePeriod).Methods("PUT")
	router.HandleFunc("/booking/unavailability/add", handler.AddPeriod).Methods("PUT")
}

func (handler *UnavailabilityHandler) AddPeriod(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "add-post")
	defer func() { span.End() }()
	var manageUnavailabilityPeriodDto dto.ManageUnavailabilityPeriodDto
	if err := json.NewDecoder(r.Body).Decode(&manageUnavailabilityPeriodDto); err != nil {
		util.HttpTraceError(err, "invalid request payload", span, handler.loki, "AddPeriod", "")
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateManageUnavailabilityPeriodDto(manageUnavailabilityPeriodDto); err != nil {
		util.HttpTraceError(err, "invalid request data", span, handler.loki, "AddPeriod", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	newUnavailabilityPeriod := dto.MapUnavailabilityPeriod(&manageUnavailabilityPeriodDto)
	newUnavailabilityPeriod.Reason = domain.OwnerSet
	if err := handler.service.AddUnavailabilityPeriod(manageUnavailabilityPeriodDto.AccommodationId, newUnavailabilityPeriod, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to add unavailability period", span, handler.loki, "AddPeriod", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Unavailability period added successfully", span, handler.loki, "AddPeriod", "")

	w.WriteHeader(http.StatusCreated)
}

func (handler *UnavailabilityHandler) DeletePeriod(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "delete-put")
	defer func() { span.End() }()
	var manageUnavailabilityPeriodDto dto.ManageUnavailabilityPeriodDto
	if err := json.NewDecoder(r.Body).Decode(&manageUnavailabilityPeriodDto); err != nil {
		util.HttpTraceError(err, "invalid request payload", span, handler.loki, "DeletePeriod", "")
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateManageUnavailabilityPeriodDto(manageUnavailabilityPeriodDto); err != nil {
		util.HttpTraceError(err, "invalid request data", span, handler.loki, "DeletePeriod", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	removedUnavailabilityPeriod := dto.MapUnavailabilityPeriod(&manageUnavailabilityPeriodDto)
	if err := handler.service.RemoveUnavailabilityPeriod(manageUnavailabilityPeriodDto.AccommodationId, removedUnavailabilityPeriod, true, span, handler.loki); err != nil {
		util.HttpTraceError(err, "failed to remove unavailability period", span, handler.loki, "DeletePeriod", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Unavailability period removed successfully", span, handler.loki, "DeletePeriod", "")

	w.WriteHeader(http.StatusCreated)
}

func (handler *UnavailabilityHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-all-get")
	defer func() { span.End() }()
	unavailabilityList, err := handler.service.GetAll(span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "failed to get all unavailability", span, handler.loki, "GetAll", "")
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
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetAll", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Unavailability successfully fetched", span, handler.loki, "GetAll", "")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *UnavailabilityHandler) GetByAccommodationId(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-accommodation-id-get")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	accommodationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "GetByAccommodationId", "")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unavailability, err := handler.service.GetByAccommodationId(accommodationId, span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "failed to get by accommodation id", span, handler.loki, "GetByAccommodationId", "")
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
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetByAccommodationId", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Unavailability from accommodation fetched successfully", span, handler.loki, "GetByAccommodationId", "")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *UnavailabilityHandler) GetByHostId(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-host-id")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	hostId := vars["id"]
	unavailabilityList, err := handler.service.GetByHostId(hostId, span, handler.loki)

	if err != nil {
		util.HttpTraceError(err, "failed to parse form data", span, handler.loki, "GetByHostId", "")
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
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetByHostId", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	util.HttpTraceInfo("Unavailability fetched by host id successfully", span, handler.loki, "GetByHostId", "")

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
