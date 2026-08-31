package user

import "time"

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"    validate:"omitempty,min=3,max=50"`
	Password    string `json:"password"    validate:"omitempty,min=8"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
}

type Response struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewResponse(user *User) Response {
	return Response{
		ID:          user.ID,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}
}
