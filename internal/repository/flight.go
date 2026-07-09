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

type FlightRepositoryInterface interface {
	Create(ctx context.Context, flight *Flight) (*Flight, error)
	GetAll(ctx context.Context) ([]Flight, error)
	GetByID(ctx context.Context, id int) (Flight, error)
	Update(ctx context.Context, flight *Flight) (*Flight, error)
	DeleteByID(ctx context.Context, id int) error
}

type Flight struct {
	ID          int       `json:"id"`
	AircraftID  int       `json:"aircraft_id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	DepartureAt time.Time `json:"departure_at"`
	ArriveAt    time.Time `json:"arrive_at"`
}

type FlightRepository struct {
	db *pgxpool.Pool
}

func NewFlightRepository(db *pgxpool.Pool) FlightRepositoryInterface {
	return &FlightRepository{db: db}
}

func (r *FlightRepository) Create(ctx context.Context, flight *Flight) (*Flight, error) {
	query := "INSERT INTO flights (aircraft_id, from_location, to_location, departure_at, arrive_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := r.db.QueryRow(
		ctx,
		query,
		flight.AircraftID,
		flight.From,
		flight.To,
		flight.DepartureAt,
		flight.ArriveAt,
	).Scan(&flight.ID)

	if err != nil {
		logger.Warn(err)
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return flight, nil
}

func (r *FlightRepository) GetAll(ctx context.Context) ([]Flight, error) {
	flights := make([]Flight, 0)
	query := "SELECT * FROM flights"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Warn(err)
		return nil, apperror.Internal("Internal Server Error", err)
	}

	defer rows.Close()

	for rows.Next() {
		var flight Flight

		err := rows.Scan(
			&flight.ID,
			&flight.AircraftID,
			&flight.From,
			&flight.To,
			&flight.DepartureAt,
			&flight.ArriveAt,
		)

		if err != nil {
			logger.Warn(err)
			return nil, apperror.Internal("Internal Server Error", err)
		}

		flights = append(flights, flight)
	}

	return flights, nil
}

func (r *FlightRepository) GetByID(ctx context.Context, id int) (Flight, error) {
	var flight Flight

	query := "SELECT id, aircraft_id, from_location, to_location, departure_at, arrive_at FROM flights WHERE id=$1"

	err := r.db.QueryRow(ctx, query, id).Scan(
		&flight.ID,
		&flight.AircraftID,
		&flight.From,
		&flight.To,
		&flight.DepartureAt,
		&flight.ArriveAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Flight{}, apperror.NotFound("Flight not found")
		}

		return Flight{}, apperror.Internal("Internal Server Error", err)
	}

	return flight, nil
}

func (r *FlightRepository) Update(ctx context.Context, flight *Flight) (*Flight, error) {
	query := `
				UPDATE flights 
				SET aircraft_id=$1, from_location=$2, to_location=$3, departure_at=$4, arrive_at=$5 
				WHERE id=$6
				RETURNING id, aircraft_id, from_location, to_location, departure_at, arrive_at`

	err := r.db.QueryRow(
		ctx,
		query,
		flight.AircraftID,
		flight.From,
		flight.To,
		flight.DepartureAt,
		flight.ArriveAt,
		flight.ID,
	).Scan(
		&flight.ID,
		&flight.AircraftID,
		&flight.From,
		&flight.To,
		&flight.DepartureAt,
		&flight.ArriveAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn(err)
			return nil, apperror.NotFound("Flight not found")
		}

		logger.Warn(err)
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return flight, nil
}

func (r *FlightRepository) DeleteByID(ctx context.Context, id int) error {
	query := "DELETE FROM flights WHERE id=$1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Warn(err)
		return apperror.Internal("Internal Server Error", err)
	}

	if result.RowsAffected() == 0 {
		logger.Warn(err)
		return apperror.NotFound("Flight not found")
	}

	return nil
}
