package repository

import (
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
	GetByID(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, user *User) (*User, error)
	DeleteByID(ctx context.Context, id int) error
}

type User struct {
	ID          int
	Username    string
	Password    string
	PhoneNumber string
	CreatedAt   time.Time
	DeletedAt   *time.Time
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

func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	users := make([]User, 0)
	query := "SELECT * FROM users"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Warn(err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.PhoneNumber, &user.CreatedAt, &user.DeletedAt)
		if err != nil {
			logger.Warn(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (User, error) {
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
			return User{}, nil
		}
		logger.Warn(err)
		return User{}, nil
	}
	return user, nil
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
			return nil, errors.New("user not found")
		}
		logger.Warn(err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id=$1"

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn(err)
			return err
		}
		logger.Warn(err)
		return err
	}
	return nil
}
