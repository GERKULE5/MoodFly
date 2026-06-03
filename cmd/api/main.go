package main

import (
	"MoodFly/pkg/database"
	"MoodFly/pkg/logger"
)

func main() {
	logger.Init()
	db, err := database.ConnectDB()
	if err != nil {
		logger.Err("Failed to connect DB")
	}
	defer db.Close()
}
