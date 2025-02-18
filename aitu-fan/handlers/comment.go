package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/middleware"
	"github.com/group-project/aitu-fan/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostID  string `json:"post_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userIDStr, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	postObjectID, err := primitive.ObjectIDFromHex(req.PostID)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByID(userIDStr)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	comment := models.Comment{
		ID:        primitive.NewObjectID(),
		PostID:    postObjectID,
		UserID:    userObjectID,
		Username:  user.Username,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	coll := database.GetCollection("comments")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = coll.InsertOne(ctx, comment)
	if err != nil {
		http.Error(w, "Error adding comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "post_id is required", http.StatusBadRequest)
		return
	}
	postObjectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	coll := database.GetCollection("comments")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"post_id": postObjectID}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		http.Error(w, "Error processing comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
