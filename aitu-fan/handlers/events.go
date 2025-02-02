package handlers

import (
	"encoding/json"
	"github.com/group-project/aitu-fan/aitu-fan/database"
	"github.com/group-project/aitu-fan/aitu-fan/models"
	"net/http"
	"time"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event.Date = time.Now()

	result := database.DB.Create(&event)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	var events []models.Event

	result := database.DB.Find(&events)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
