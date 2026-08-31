package user

import (
	apperror "MoodFly/pkg/error"
	"MoodFly/pkg/logger"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo     Repository
	validate *validator.Validate
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *Service) Create(ctx context.Context, data *CreateUserRequest) (*User, error) {
	err := s.validate.Struct(data)
	if err != nil {
		logger.Warn(err)
		return nil, apperror.BadRequest("Failed to create user")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("Internal Server Error", err)
	}

	user := &User{
		Username:    data.Username,
		Password:    string(hashed),
		PhoneNumber: data.PhoneNumber,
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) GetAll(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) Update(ctx context.Context, id int, data *UpdateUserRequest) (*User, error) {
	if err := s.validate.Struct(data); err != nil {
		logger.Warn(err)
		// ПРОБЛЕМА: сырой validator error не является AppError, поэтому
		// HandleError превратит ошибку клиента в HTTP 500.
		// РЕШЕНИЕ: вернуть apperror.BadRequest("Invalid user data").
		return nil, apperror.BadRequest("Invalid user data")
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if data.Username != "" {
		existing.Username = data.Username
	}
	if data.PhoneNumber != "" {
		existing.PhoneNumber = data.PhoneNumber
	}
	if data.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Err("failed to hash password: ", err)
			return nil, apperror.Internal("Internal Server Error", err)
		}
		existing.Password = string(hashed)
	}

	return s.repo.Update(ctx, existing)
}

func (s *Service) DeleteByID(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
