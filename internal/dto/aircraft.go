package dto

import "time"

type CreateAircraftDTO struct {
	WINNum     string    `json:"win_num" validate:"required,min=12,max=12"`
	Model      string    `json:"model" validate:"required,max=100"`
	Capacity   int       `json:"capacity" validate:"required"`
	CarryingKg int       `json:"carrying_kg" validate:"required,min=0"`
	FlightTime int       `json:"flight_time" validate:"required,min=0"`
	ReleasedAt time.Time `json:"released_at" validate:"required"`
	LicensedAt time.Time `json:"licensed_at" validate:"required"`
}

type UpdateAircraftDTO struct {
	WINNum     *string    `json:"win_num" validate:"omitempty,min=12,max=12"`
	Model      *string    `json:"model" validate:"omitempty,max=100"`
	Capacity   *int       `json:"capacity" validate:"omitempty"`
	CarryingKg *int       `json:"carrying_kg" validate:"omitempty,min=0"`
	FlightTime *int       `json:"flight_time" validate:"omitempty,min=0"`
	ReleasedAt *time.Time `json:"released_at" validate:"omitempty"`
	LicensedAt *time.Time `json:"licensed_at" validate:"omitempty"`
}
