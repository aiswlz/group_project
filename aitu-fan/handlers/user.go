package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/group-project/aitu-fan/aitu-fan/database"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Unable to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	err := database.DeleteUserByID(userID)
	if err != nil {
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintf(w, "User with ID %s deleted successfully", userID)
}
