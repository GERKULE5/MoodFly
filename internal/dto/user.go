package dto

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"    validate:"omitempty,min=3,max=50"`
	Password    string `json:"password"    validate:"omitempty,min=6"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,e164"`
}
