package handler

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/service"
	apperror "MoodFly/pkg/error"
	"encoding/json"
	"net/http"
	"strconv"
)

type FlightHandler struct {
	service service.FlightServiceInterface
}

func NewFlightHandler(service service.FlightServiceInterface) *FlightHandler {
	return &FlightHandler{
		service: service,
	}
}

func (handler *FlightHandler) CreateFlight(w http.ResponseWriter, r *http.Request) {
	var flight dto.CreateFlightDTO

	err := json.NewDecoder(r.Body).Decode(&flight)
	if err != nil {
		handleError(w, apperror.BadRequest("Bad Request"))
		return
	}

	result, err := handler.service.Create(r.Context(), &flight)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (handler *FlightHandler) GetAllFlights(w http.ResponseWriter, r *http.Request) {
	flights, err := handler.service.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flights)
}

func (handler *FlightHandler) GetFlightByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, apperror.BadRequest("Invalid ID"))
		return
	}

	flight, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flight)
}

func (handler *FlightHandler) UpdateFlight(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, apperror.BadRequest("Invalid ID"))
		return
	}

	var flight dto.UpdateFlightDTO

	err = json.NewDecoder(r.Body).Decode(&flight)
	if err != nil {
		handleError(w, apperror.BadRequest("Bad Request"))
		return
	}

	result, err := handler.service.Update(r.Context(), id, &flight)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (handler *FlightHandler) DeleteFlight(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, apperror.BadRequest("Invalid ID"))
		return
	}

	err = handler.service.Delete(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
