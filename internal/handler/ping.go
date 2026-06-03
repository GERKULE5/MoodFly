package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type PingResponse struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	response := &PingResponse{
		Status: "OK",
		Time:   time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
