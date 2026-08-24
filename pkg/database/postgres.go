package database

import (
	"MoodFly/pkg/logger"
	"MoodFly/pkg/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	config := Config{
		Host:     utils.GetEnv("POSTGRES_HOST", "localhost"),
		Port:     utils.GetEnv("POSTGRES_PORT", "5432"),
		User:     utils.GetEnv("POSTGRES_USER", "postgres"),
		Password: utils.GetEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:   utils.GetEnv("POSTGRES_DB", "crud_db"),
		SSLMode:  utils.GetEnv("POSTGRES_SSLMODE", "disable"),
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Err("Failed connect to DB: ", err)
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		logger.Err("Unable to Ping DB: ", err)
		return nil, err
	}

	logger.Info("Successfully connected to DB")
	return db, nil
}
