package main

import (
	"MoodFly/internal/handler"
	"MoodFly/internal/middleware"
	"MoodFly/internal/repository"
	"MoodFly/internal/service"
	"MoodFly/pkg/database"
	"MoodFly/pkg/logger"
	"MoodFly/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func protected(h http.HandlerFunc) http.Handler {
	return middleware.AuthMiddleware(h)
}

func main() {
	_ = godotenv.Load()
	logger.Init()

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db, err := database.ConnectDB(ctx)
	if err != nil {
		logger.Err("Failed to connect DB")
		return
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + utils.GetEnv("REDIS_PORT", "6379"),
		Password: utils.GetEnv("REDIS_PASSWORD", "root"),
		DB:       1,
	})
	defer rdb.Close()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Err("Could not connect to Redis: %v", err)
		return
	}
	logger.Info("Connected to Redis successfully:", pong)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handler.Ping)

	// USERS
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.Handle("GET /users", protected(userHandler.GetAllUsers))
	mux.Handle("GET /users/{id}", protected(userHandler.GetUserByID))
	mux.Handle("PUT /users/{id}", protected(userHandler.UpdateUser))
	mux.Handle("DELETE /users/{id}", protected(userHandler.DeleteUserByID))

	// AUTH
	authService := service.NewAuthService(rdb)
	authHandler := handler.NewAuthHandler(authService, userService)

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /refresh", authHandler.Refresh)
	mux.Handle("POST /logout", protected(authHandler.Logout))

	// AIRCRAFTS
	aircraftRepo := repository.NewAircraftRepository(db)
	aircraftService := service.NewAircraftService(aircraftRepo)
	aircraftHandler := handler.NewAircraftHandler(aircraftService)

	mux.Handle("POST /aircrafts", protected(aircraftHandler.CreateAircraft))
	mux.Handle("GET /aircrafts", protected(aircraftHandler.GetAllAircrafts))
	mux.Handle("GET /aircrafts/{id}", protected(aircraftHandler.GetAircraftByID))
	mux.Handle("PUT /aircrafts/{id}", protected(aircraftHandler.UpdateAircraft))
	mux.Handle("DELETE /aircrafts/{id}", protected(aircraftHandler.DeleteAircraft))

	// FLIGHTS
	flightRepo := repository.NewFlightRepository(db)
	flightService := service.NewFlightService(flightRepo, aircraftService)
	flightHandler := handler.NewFlightHandler(flightService)

	mux.Handle("POST /flights", protected(flightHandler.CreateFlight))
	mux.Handle("GET /flights", protected(flightHandler.GetAllFlights))
	mux.Handle("GET /flights/{id}", protected(flightHandler.GetFlightByID))
	mux.Handle("PUT /flights/{id}", protected(flightHandler.UpdateFlight))
	mux.Handle("DELETE /flights/{id}", protected(flightHandler.DeleteFlight))

	logger.Info("Server started on 8080")

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Err("Server error: ", err)
		return
	}
}
