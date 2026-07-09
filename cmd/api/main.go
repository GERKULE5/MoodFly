package main

import (
	"MoodFly/internal/handler"
	"MoodFly/internal/repository"
	"MoodFly/internal/service"
	"MoodFly/pkg/database"
	"MoodFly/pkg/logger"
	"net/http"
)

func main() {
	logger.Init()
	db, err := database.ConnectDB()
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
	userHandler := handler.NewUserNandler(userService)

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

	logger.Info("Server started on 3000")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		logger.Err("Server error: ", err)
		return
	}
}
