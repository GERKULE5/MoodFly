package aircraft

import (
	apperror "MoodFly/pkg/error"
	"MoodFly/pkg/logger"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
)

type AircraftServiceInterface interface {
	Create(ctx context.Context, data *CreateAircraftDto) (*Aircraft, error)
	GetAll(ctx context.Context) ([]Aircraft, error)
	GetByID(ctx context.Context, id int) (Aircraft, error)
	Update(ctx context.Context, id int, data *UpdateAircraftDto) (*Aircraft, error)
	Delete(ctx context.Context, id int) error
	IsAvailable(ctx context.Context, aircraftID int, newDepartureAt time.Time, newArriveAt time.Time) (bool, error)
}

type AircraftService struct {
	repository AircraftRepositoryInterface
	validate   *validator.Validate
}

func NewAircraftService(repository AircraftRepositoryInterface) AircraftServiceInterface {
	return &AircraftService{
		repository: repository,
		validate:   validator.New(),
	}
}

func (service *AircraftService) Create(ctx context.Context, data *CreateAircraftDto) (*Aircraft, error) {
	err := service.validate.Struct(data)
	if err != nil {
		logger.Warn(err)
		return nil, apperror.BadRequest("Bad request")
	}

	aircraft := &Aircraft{
		WINNum:     data.WINNum,
		Model:      data.Model,
		Capacity:   data.Capacity,
		CarryingKg: data.CarryingKg,
		FlightTime: data.FlightTime,
		ReleasedAt: data.ReleasedAt,
		LicensedAt: data.LicensedAt,
	}

	return service.repository.Create(ctx, aircraft)
}

func (service *AircraftService) GetAll(ctx context.Context) ([]Aircraft, error) {
	return service.repository.GetAll(ctx)
}

func (service *AircraftService) GetByID(ctx context.Context, id int) (Aircraft, error) {
	return service.repository.GetByID(ctx, id)
}

func (service *AircraftService) Update(ctx context.Context, id int, data *UpdateAircraftDto) (*Aircraft, error) {
	err := service.validate.Struct(data)

	if err != nil {
		logger.Warn(err)
		return nil, apperror.BadRequest("Bad Request")
	}

	existing, err := service.repository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if data.WINNum != nil {
		existing.WINNum = *data.WINNum
	}

	if data.Model != nil {
		existing.Model = *data.Model
	}

	if data.Capacity != nil {
		existing.Capacity = *data.Capacity
	}

	if data.CarryingKg != nil {
		existing.CarryingKg = *data.CarryingKg
	}

	if data.FlightTime != nil {
		existing.FlightTime = *data.FlightTime
	}

	if data.ReleasedAt != nil {
		existing.ReleasedAt = *data.ReleasedAt
	}

	if data.LicensedAt != nil {
		existing.LicensedAt = *data.LicensedAt
	}

	return service.repository.Update(ctx, &existing)
}

func (service *AircraftService) Delete(ctx context.Context, id int) error {
	return service.repository.DeleteByID(ctx, id)
}

func (service *AircraftService) IsAvailable(ctx context.Context, aircraftID int, newDepartureAt time.Time, newArriveAt time.Time) (bool, error) {
	return service.repository.IsAvailable(ctx, aircraftID, newDepartureAt, newArriveAt)
}
