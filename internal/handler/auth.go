package handler

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/service"
	apperror "MoodFly/pkg/error"
	auth_jwt "MoodFly/pkg/jwt"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	service service.AuthServiceInterface
	users   service.UserServiceInterface
}

func NewAuthHandler(service service.AuthServiceInterface, users service.UserServiceInterface) *AuthHandler {
	return &AuthHandler{
		service: service,
		users:   users,
	}
}

func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, apperror.BadRequest("invalid request body"))
		return
	}

	user, err := handler.users.Create(r.Context(), &dto.CreateUserRequest{
		Username:    request.Username,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
	})

	if err != nil {
		handleError(w, err)
		return
	}

	access, err := auth_jwt.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		handleError(w, err)
	}

	refresh, err := handler.service.CreateToken(r.Context(), user.ID, user.Username)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, apperror.BadRequest("invalid request body"))
		return
	}

	user, err := handler.users.GetByUsername(r.Context(), request.Username)
	if err != nil {
		handleError(w, apperror.Unauthorized("invalid credentials"))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		handleError(w, apperror.Unauthorized("invalid credentials"))
		return
	}

	access, err := auth_jwt.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		handleError(w, apperror.Internal("could not generate token", err))
		return
	}

	refresh, err := handler.service.CreateToken(r.Context(), user.ID, user.Username)
	if err != nil {
		handleError(w, apperror.Internal("could not create session", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse{AccessToken: access, RefreshToken: refresh})
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var request refreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, apperror.BadRequest("invalid request body"))
		return
	}

	userID, err := handler.service.Validate(r.Context(), request.RefreshToken)
	if err != nil {
		handleError(w, apperror.Unauthorized("invalid or expired refresh token"))
		return
	}

	if err := handler.service.Revoke(r.Context(), request.RefreshToken); err != nil {
		handleError(w, apperror.Internal("could not revoke old token", err))
		return
	}

	access, err := auth_jwt.GenerateAccessToken(userID, "")
	if err != nil {
		handleError(w, apperror.Internal("could not generate token", err))
		return
	}

	newRefresh, err := handler.service.CreateToken(r.Context(), userID, "")
	if err != nil {
		handleError(w, apperror.Internal("could not create session", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse{AccessToken: access, RefreshToken: newRefresh})
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var request refreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handleError(w, apperror.BadRequest("invalid request body"))
		return
	}

	if err := handler.service.Revoke(r.Context(), request.RefreshToken); err != nil {
		handleError(w, apperror.Internal("could not revoke token", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
