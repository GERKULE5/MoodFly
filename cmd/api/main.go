package main

import (
	"MoodFly/pkg/database"
	"fmt"
)

func main() {
	fmt.Println("dfdf")

	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect DB")
	}
	defer db.Close()
}
