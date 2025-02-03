package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/group-project/aitu-fan/aitu-fan/database"
	"github.com/group-project/aitu-fan/aitu-fan/models"
	"net/http"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	json.NewDecoder(r.Body).Decode(&comment)

	database.DB.Create(&comment)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	var comments []models.Comment
	database.DB.Find(&comments)
	json.NewEncoder(w).Encode(comments)
}

func GetCommentsByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var comments []models.Comment

	if postID, ok := params["post_id"]; ok {
		database.DB.Where("post_id = ?", postID).Find(&comments)
	} else if eventID, ok := params["event_id"]; ok {
		database.DB.Where("event_id = ?", eventID).Find(&comments)
	}

	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var comment models.Comment

	if err := database.DB.First(&comment, params["id"]).Error; err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	json.NewDecoder(r.Body).Decode(&comment)
	database.DB.Save(&comment)
	json.NewEncoder(w).Encode(comment)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	database.DB.Delete(&models.Comment{}, params["id"])
	w.WriteHeader(http.StatusNoContent)
}
