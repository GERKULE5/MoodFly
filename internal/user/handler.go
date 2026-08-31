package user

import (
	apperror "MoodFly/pkg/error"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.service.Create(r.Context(), &request)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response := NewResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (handler *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.service.GetAll(r.Context())
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response := make([]Response, 0, len(users))

	for i := range users {
		response = append(response, NewResponse(&users[i]))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (handler *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response := NewResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (handler *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var request UpdateUserRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.service.Update(r.Context(), id, &request)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}

	response := NewResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (handler *Handler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = handler.service.DeleteByID(r.Context(), id)
	if err != nil {
		apperror.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
