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

func LikePost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PostID string `json:"post_id"`
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

	coll := database.GetCollection("likes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"post_id": postObjectID, "user_id": userObjectID}
	var existing models.Like
	err = coll.FindOne(ctx, filter).Decode(&existing)
	if err == nil {

		_, err = coll.DeleteOne(ctx, filter)
		if err != nil {
			http.Error(w, "Error unliking post", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unliked"})
		return
	}

	like := models.Like{
		ID:     primitive.NewObjectID(),
		PostID: postObjectID,
		UserID: userObjectID,
	}
	_, err = coll.InsertOne(ctx, like)
	if err != nil {
		http.Error(w, "Error liking post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(like)
}

func GetLikes(w http.ResponseWriter, r *http.Request) {
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

	coll := database.GetCollection("likes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"post_id": postObjectID}
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		http.Error(w, "Error counting likes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"likes": count})
}
