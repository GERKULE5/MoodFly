package dto

import "time"

type CreateFlightDTO struct {
	AircraftID  int       `json:"aircraft_id" validate:"required"`
	From        string    `json:"from" validate:"required,max=100"`
	To          string    `json:"to" validate:"required,max=100"`
	DepartureAt time.Time `json:"departure_at" validate:"required"`
	ArriveAt    time.Time `json:"arrive_at" validate:"required"`
}

type UpdateFlightDTO struct {
	AircraftID  *int       `json:"aircraft_id" validate:"omitempty"`
	From        *string    `json:"from" validate:"omitempty,max=100"`
	To          *string    `json:"to" validate:"omitempty,max=100"`
	DepartureAt *time.Time `json:"departure_at" validate:"omitempty"`
	ArriveAt    *time.Time `json:"arrive_at" validate:"omitempty"`
}
