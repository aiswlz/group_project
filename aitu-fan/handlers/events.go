package handlers

import (
	"encoding/json"
	"github.com/group-project/aitu-fan/aitu-fan/database"
	"github.com/group-project/aitu-fan/aitu-fan/models"
	"net/http"
	"time"
)

func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := database.GetAllEvents()
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event.Date = time.Now()

	if err := database.CreateEvent(event); err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event created successfully"})
}
