package repository

import (
	"MoodFly/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) (*User, error)
}

type User struct {
	ID          int
	Username    string
	Password    string
	PhoneNumber string
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) (*User, error) {
	query := "INSERT INTO users (username, password, phone_number) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(ctx, query, user.Username, user.Password, user.PhoneNumber).Scan(&user.ID)
	if err != nil {
		logger.Warn(err)
		return nil, err
	}

	return user, nil
}
