package service

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/repository"
	apperror "MoodFly/pkg/error"
	"MoodFly/pkg/logger"
	"context"

	"github.com/go-playground/validator/v10"
)

type AircraftServiceInterface interface {
	Create(ctx context.Context, data *dto.CreateAircraftDTO) (*repository.Aircraft, error)
	GetAll(ctx context.Context) ([]repository.Aircraft, error)
	GetByID(ctx context.Context, id int) (repository.Aircraft, error)
	Update(ctx context.Context, id int, data *dto.UpdateAircraftDTO) (*repository.Aircraft, error)
	Delete(ctx context.Context, id int) error
}

type AircraftService struct {
	repository repository.AircraftRepositoryInterface
	validate   *validator.Validate
}

func NewAircraftService(repository repository.AircraftRepositoryInterface) AircraftServiceInterface {
	return &AircraftService{
		repository: repository,
		validate:   validator.New(),
	}
}

func (service *AircraftService) Create(ctx context.Context, data *dto.CreateAircraftDTO) (*repository.Aircraft, error) {
	err := service.validate.Struct(data)
	if err != nil {
		logger.Warn(err)
		return nil, apperror.BadRequest("Bad request")
	}

	aircraft := &repository.Aircraft{
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

func (service *AircraftService) GetAll(ctx context.Context) ([]repository.Aircraft, error) {
	return service.repository.GetAll(ctx)
}

func (service *AircraftService) GetByID(ctx context.Context, id int) (repository.Aircraft, error) {
	return service.repository.GetByID(ctx, id)
}

func (service *AircraftService) Update(ctx context.Context, id int, data *dto.UpdateAircraftDTO) (*repository.Aircraft, error) {
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
