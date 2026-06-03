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
