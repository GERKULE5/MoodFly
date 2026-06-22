package handler

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/service"
	apperror "MoodFly/pkg/error"
	"encoding/json"
	"net/http"
	"strconv"
)

type AircraftHandler struct {
	service service.AircraftServiceInterface
}

func NewAircraftHandler(service service.AircraftServiceInterface) *AircraftHandler {
	return &AircraftHandler{
		service: service,
	}
}

func (handler *AircraftHandler) CreateAircraft(w http.ResponseWriter, r *http.Request) {
	var aircraft dto.CreateAircraftDTO

	err := json.NewDecoder(r.Body).Decode(&aircraft)
	if err != nil {
		handleError(w, apperror.BadRequest("Bad Requestdf"))
		return
	}

	result, err := handler.service.Create(r.Context(), &aircraft)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (handler *AircraftHandler) GetAllAircrafts(w http.ResponseWriter, r *http.Request) {
	aicrafts, err := handler.service.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(aicrafts)
}

func (handler *AircraftHandler) GetAircraftByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, apperror.BadRequest("Invalid ID"))
		return
	}

	aircraft, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(aircraft)
}

func (handler *AircraftHandler) UpdateAircraft(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, apperror.BadRequest("Invalid ID"))
		return
	}

	var aircraft dto.UpdateAircraftDTO

	err = json.NewDecoder(r.Body).Decode(&aircraft)
	if err != nil {
		handleError(w, apperror.BadRequest("Bad Request"))
		return
	}

	result, err := handler.service.Update(r.Context(), id, &aircraft)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (handler *AircraftHandler) DeleteAircraft(w http.ResponseWriter, r *http.Request) {
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
