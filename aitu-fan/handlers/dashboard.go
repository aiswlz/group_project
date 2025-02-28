package handlers

import (
	"context"
	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/middleware"
	"github.com/group-project/aitu-fan/models"
	"html/template"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	postsCollection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sortOrder := -1
	if r.URL.Query().Get("sort") == "old" {
		sortOrder = 1 // сортировка от старых к новым
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: sortOrder}})

	cursor, err := postsCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		http.Error(w, "Error processing post data", http.StatusInternalServerError)
		return
	}

	userRole, err := middleware.GetUserRole(r)
	if err != nil {
		userRole = "user"
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Posts":    posts,
		"UserRole": userRole,
	})
}
