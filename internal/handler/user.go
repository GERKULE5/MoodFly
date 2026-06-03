package handler

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserNandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.service.Create(r.Context(), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}
