package main

import (
	"MoodFly/internal/handler"
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

	logger.Info("Server started on 3000")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		logger.Err("Server error: ", err)
		return
	}
}
