package service

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/repository"
	apperror "MoodFly/pkg/error"
	"MoodFly/pkg/logger"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

type FlightServiceInterface interface {
	Create(ctx context.Context, data *dto.CreateFlightDTO) (*repository.Flight, error)
	GetAll(ctx context.Context) ([]repository.Flight, error)
	GetByID(ctx context.Context, id int) (repository.Flight, error)
	Update(ctx context.Context, id int, data *dto.UpdateFlightDTO) (*repository.Flight, error)
	Delete(ctx context.Context, id int) error
}

type FlightService struct {
	repository      repository.FlightRepositoryInterface
	aircraftService AircraftServiceInterface
	validate        *validator.Validate
}

func NewFlightService(repository repository.FlightRepositoryInterface, aircraftService AircraftServiceInterface) FlightServiceInterface {
	return &FlightService{
		repository:      repository,
		aircraftService: aircraftService,
		validate:        validator.New(),
	}
}

func (service *FlightService) Create(ctx context.Context, data *dto.CreateFlightDTO) (*repository.Flight, error) {
	err := service.validate.Struct(data)
	if err != nil {
		logger.Warn(err)
		return nil, apperror.BadRequest("Bad request")
	}

	correctDate := checkDate(data.DepartureAt, data.ArriveAt)
	if correctDate == false {
		logger.Warn("Incorrect date")
		return nil, apperror.BadRequest("Departure time must be before arrival time")
	}

	available, err := service.aircraftService.IsAvailable(ctx, data.AircraftID, data.DepartureAt, data.ArriveAt)

	if available == false || err != nil {
		if err != nil {
			logger.Warn(err)
		}
		logger.Warn("Aircraft is busy")
		return nil, apperror.Conflict("Aircraft is busy", err)
	}

	flight := &repository.Flight{
		AircraftID:  data.AircraftID,
		From:        data.From,
		To:          data.To,
		DepartureAt: data.DepartureAt,
		ArriveAt:    data.ArriveAt,
	}

	return service.repository.Create(ctx, flight)
}

func (service *FlightService) GetAll(ctx context.Context) ([]repository.Flight, error) {
	return service.repository.GetAll(ctx)
}

func (service *FlightService) GetByID(ctx context.Context, id int) (repository.Flight, error) {
	return service.repository.GetByID(ctx, id)
}

func (service *FlightService) Update(ctx context.Context, id int, data *dto.UpdateFlightDTO) (*repository.Flight, error) {
	err := service.validate.Struct(data)

	if err != nil {
		logger.Warn(err)
		return nil, apperror.BadRequest("Bad Request")
	}

	existing, err := service.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if data.AircraftID != nil {
		existing.AircraftID = *data.AircraftID
	}

	if data.From != nil {
		existing.From = *data.From
	}

	if data.To != nil {
		existing.To = *data.To
	}

	if data.DepartureAt != nil {
		existing.DepartureAt = *data.DepartureAt
	}

	if data.ArriveAt != nil {
		existing.ArriveAt = *data.ArriveAt
	}

	return service.repository.Update(ctx, &existing)
}

func (service *FlightService) Delete(ctx context.Context, id int) error {
	return service.repository.DeleteByID(ctx, id)
}

func checkDate(departureAt, arriveAt time.Time) bool {
	if departureAt.After(arriveAt) || departureAt.Equal(arriveAt) {
		return false
	}

	return true
}
