package repository

import (
	apperror "MoodFly/pkg/error"
	"MoodFly/pkg/logger"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	DeleteByID(ctx context.Context, id int) error
}

type User struct {
	ID          int        `json:"id"`
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	PhoneNumber string     `json:"phone_number"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
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
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	users := make([]User, 0)
	query := "SELECT * FROM users WHERE deleted_at IS NULL"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Warn(err)
		return users, apperror.Internal("Internal Server Error", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.PhoneNumber, &user.CreatedAt, &user.DeletedAt)
		if err != nil {
			logger.Warn(err)
			return nil, apperror.Internal("Internal Server Error", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	var user User

	query := "SELECT id, username, password, phone_number, created_at, deleted_at FROM users WHERE id=$1"

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn(err)
			return &User{}, apperror.NotFound("User not found")
		}
		logger.Warn(err)
		return &User{}, apperror.Internal("Internal Server Error", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	query := "SELECT id, username, password, phone_number, created_at, deleted_at FROM users WHERE username=$1"

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn(err)
			return &User{}, apperror.NotFound("User not found")
		}
		logger.Warn(err)
		return &User{}, apperror.Internal("Internal Server Error", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) (*User, error) {
	query := `
        UPDATE users
        SET username = $1, password = $2, phone_number = $3
        WHERE id = $4
        RETURNING id, username, password, phone_number, created_at, deleted_at`

	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.Password,
		user.PhoneNumber,
		user.ID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("User not found")
		}
		logger.Warn(err)
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return user, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id int) error {
	query := "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Warn(err)
		return apperror.Internal("Internal Server Error", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NotFound("User not found")
	}

	return nil
}
