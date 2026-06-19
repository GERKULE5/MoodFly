package handler

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/service"
	apperror "MoodFly/pkg/error"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func handleError(w http.ResponseWriter, err error) {
	code, msg := apperror.ToHTTP(err)
	http.Error(w, msg, code)
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
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

}

func (handler *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.service.GetAll(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (handler *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := handler.service.GetByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var request dto.UpdateUserRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := handler.service.Update(r.Context(), id, &request)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (handler *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = handler.service.DeleteByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
