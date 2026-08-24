package dto

type RegisterDTO struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type LoginDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}
