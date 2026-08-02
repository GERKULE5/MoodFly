package repository

import (
	apperror "MoodFly/pkg/error"
	"MoodFly/pkg/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AircraftRepositoryInterface interface {
	Create(ctx context.Context, aircraft *Aircraft) (*Aircraft, error)
	GetAll(ctx context.Context) ([]Aircraft, error)
	GetByID(ctx context.Context, id int) (Aircraft, error)
	Update(ctx context.Context, aircraft *Aircraft) (*Aircraft, error)
	DeleteByID(ctx context.Context, id int) error
	IsAvailable(ctx context.Context, aircraftID int, newDepartureAt time.Time, newArriveAt time.Time) (bool, error)
}

type Aircraft struct {
	ID         int       `json:"id"`
	WINNum     string    `json:"win_num"`
	Model      string    `json:"model"`
	Capacity   int       `json:"capacity"`
	CarryingKg int       `json:"carrying_kg"`
	FlightTime int       `json:"flight_time"`
	ReleasedAt time.Time `json:"released_at"`
	LicensedAt time.Time `json:"licensed_at"`
}

type AircraftRepository struct {
	db *pgxpool.Pool
}

func NewAircraftRepository(db *pgxpool.Pool) AircraftRepositoryInterface {
	return &AircraftRepository{db: db}
}

func (r *AircraftRepository) Create(ctx context.Context, aircraft *Aircraft) (*Aircraft, error) {
	query := "INSERT INTO aircrafts (win_num, model, capacity, carrying_kg, flight_time, released_at, licensed_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	err := r.db.QueryRow(
		ctx,
		query,
		aircraft.WINNum,
		aircraft.Model,
		aircraft.Capacity,
		aircraft.CarryingKg,
		aircraft.FlightTime,
		aircraft.ReleasedAt,
		aircraft.LicensedAt,
	).Scan(&aircraft.ID)

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			logger.Warn(err)
			message := fmt.Sprintf("aircraft with this %s already exists", pgErr.ConstraintName)
			return nil, apperror.Conflict(message, err)
		}
		logger.Warn(err)
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return aircraft, nil
}

func (r *AircraftRepository) GetAll(ctx context.Context) ([]Aircraft, error) {
	aircrafts := make([]Aircraft, 0)
	query := "SELECT * FROM aircrafts"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		logger.Warn(err)
		return nil, apperror.Internal("Internal", err)
	}

	defer rows.Close()

	for rows.Next() {
		var aircraft Aircraft

		err := rows.Scan(
			&aircraft.ID,
			&aircraft.WINNum,
			&aircraft.Model,
			&aircraft.Capacity,
			&aircraft.CarryingKg,
			&aircraft.FlightTime,
			&aircraft.ReleasedAt,
			&aircraft.LicensedAt,
		)

		if err != nil {
			logger.Warn(err)
			return nil, apperror.Internal("Internal Server Error", err)
		}

		aircrafts = append(aircrafts, aircraft)
	}

	return aircrafts, nil
}

func (r *AircraftRepository) GetByID(ctx context.Context, id int) (Aircraft, error) {
	var aircraft Aircraft

	query := "SELECT id , win_num, model, capacity, carrying_kg, flight_time, released_at, licensed_at FROM aircrafts WHERE id=$1"

	err := r.db.QueryRow(ctx, query, id).Scan(
		&aircraft.ID,
		&aircraft.WINNum,
		&aircraft.Model,
		&aircraft.Capacity,
		&aircraft.CarryingKg,
		&aircraft.FlightTime,
		&aircraft.ReleasedAt,
		&aircraft.LicensedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Aircraft{}, apperror.NotFound("Not found aircrafts")
		}

		return Aircraft{}, apperror.Internal("Internal server Error", err)
	}

	return aircraft, nil
}

func (r *AircraftRepository) Update(ctx context.Context, aircraft *Aircraft) (*Aircraft, error) {
	query := `
				UPDATE aircrafts 
				SET win_num=$1, model=$2, capacity=$3, carrying_kg=$4, flight_time=$5, released_at=$6, licensed_at=$7 
				WHERE id=$8
				RETURNING id , win_num, model, capacity, carrying_kg, flight_time, released_at, licensed_at`

	err := r.db.QueryRow(
		ctx,
		query,
		aircraft.WINNum,
		aircraft.Model,
		aircraft.Capacity,
		aircraft.CarryingKg,
		aircraft.FlightTime,
		aircraft.ReleasedAt,
		aircraft.LicensedAt,
		aircraft.ID,
	).Scan(
		&aircraft.ID,
		&aircraft.WINNum,
		&aircraft.Model,
		&aircraft.Capacity,
		&aircraft.CarryingKg,
		&aircraft.FlightTime,
		&aircraft.ReleasedAt,
		&aircraft.LicensedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn(err)
			return nil, apperror.NotFound("Aircraft not found")
		}

		logger.Warn(err)
		return nil, apperror.Internal("Internal Server Error", err)
	}

	return aircraft, nil
}

func (r *AircraftRepository) DeleteByID(ctx context.Context, id int) error {
	query := "DELETE FROM aircrafts WHERE id=$1"

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Warn(err)
		return apperror.Internal("Internal Server Error", err)
	}

	if result.RowsAffected() == 0 {
		logger.Warn(err)
		return apperror.NotFound("Aircraft not found")
	}
	return nil
}

// False - busy
// True  - free
func (r *AircraftRepository) IsAvailable(ctx context.Context, aircraftID int, newDepartureAt time.Time, newArriveAt time.Time) (bool, error) {
	query := "SELECT NOT EXISTS (SELECT 1 FROM flights WHERE aircraft_id = $1 AND departure_at < $3 AND arrive_at > $2)"

	var available bool
	err := r.db.QueryRow(ctx, query, aircraftID, newDepartureAt, newArriveAt).Scan(&available)
	if err != nil {
		logger.Warn(err)
		return false, apperror.Internal("Internal Server Error", err)
	}

	return available, nil
}
