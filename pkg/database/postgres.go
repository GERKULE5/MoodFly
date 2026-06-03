package database

import (
	"MoodFly/pkg/logger"
	"context"
	"fmt"
	"os"

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

func ConnectDB() (*pgxpool.Pool, error) {
	config := Config{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnv("POSTGRES_DB", "crud_db"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Err("Failed connect to DB: ", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		logger.Err("Unable to Ping DB: ", err)
	}

	logger.Info("Successfully connected to DB")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}
