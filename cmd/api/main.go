package main

import (
	"MoodFly/internal/handler"
	"MoodFly/internal/repository"
	"MoodFly/internal/service"
	"MoodFly/pkg/database"
	"MoodFly/pkg/logger"
	"context"
	"net/http"
	"time"
)

func main() {
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", handler.Ping)

	// USERS
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserByID)

	// AIRCRAFTS
	aircraftRepo := repository.NewAircraftRepository(db)
	aircraftService := service.NewAircraftService(aircraftRepo)
	aircraftHandler := handler.NewAircraftHandler(aircraftService)

	mux.HandleFunc("POST /aircrafts", aircraftHandler.CreateAircraft)
	mux.HandleFunc("GET /aircrafts", aircraftHandler.GetAllAircrafts)
	mux.HandleFunc("GET /aircrafts/{id}", aircraftHandler.GetAircraftByID)
	mux.HandleFunc("PUT /aircrafts/{id}", aircraftHandler.UpdateAircraft)
	mux.HandleFunc("DELETE /aircrafts/{id}", aircraftHandler.DeleteAircraft)

	// FLIGHTS
	flightRepo := repository.NewFlightRepository(db)
	flightService := service.NewFlightService(flightRepo, aircraftService)
	flightHandler := handler.NewFlightHandler(flightService)

	mux.HandleFunc("POST /flights", flightHandler.CreateFlight)
	mux.HandleFunc("GET /flights", flightHandler.GetAllFlights)
	mux.HandleFunc("GET /flights/{id}", flightHandler.GetFlightByID)
	mux.HandleFunc("PUT /flights/{id}", flightHandler.UpdateFlight)
	mux.HandleFunc("DELETE /flights/{id}", flightHandler.DeleteFlight)

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
