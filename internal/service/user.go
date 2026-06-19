package service

import (
	"MoodFly/internal/dto"
	"MoodFly/internal/repository"
	"MoodFly/pkg/logger"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Create(ctx context.Context, data *dto.CreateUserRequest) (*repository.User, error)
	GetAll(ctx context.Context) ([]repository.User, error)
	GetByID(ctx context.Context, id int) (repository.User, error)
	Update(ctx context.Context, id int, data *dto.UpdateUserRequest) (*repository.User, error)
	DeleteByID(ctx context.Context, id int) error
}

type UserService struct {
	repo     repository.UserRepositoryInterface
	validate *validator.Validate
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *UserService) Create(ctx context.Context, data *dto.CreateUserRequest) (*repository.User, error) {
	err := s.validate.Struct(data)
	if err != nil {
		logger.Warn(err)
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &repository.User{
		Username:    data.Username,
		Password:    string(hashed),
		PhoneNumber: data.PhoneNumber,
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) GetAll(ctx context.Context) ([]repository.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id int) (repository.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id int, data *dto.UpdateUserRequest) (*repository.User, error) {
	if err := s.validate.Struct(data); err != nil {
		logger.Warn(err)
		return nil, err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if existing.ID == 0 {
		return nil, fmt.Errorf("user with id %d not found", id)
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
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		existing.Password = string(hashed)
	}

	return s.repo.Update(ctx, &existing)
}

func (s *UserService) DeleteByID(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}
