package handlers

import (
	"context"
	"encoding/json"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gorilla/mux"
	"github.com/group-project/aitu-fan/config"
	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/middleware"
	"github.com/group-project/aitu-fan/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
	"time"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetMyPosts(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "Failed to get userID", http.StatusUnauthorized)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": objectID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/myposts.html"))
	if err := tmpl.Execute(w, posts); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var post models.Post
	var imageURL string

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")

	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	post.UserID = objectID

	file, _, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			imageURL = ""
		} else {
			http.Error(w, "Error processing file", http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
		uploadResult, err := config.Cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{})
		if err != nil {
			http.Error(w, "Error uploading to Cloudinary", http.StatusInternalServerError)
			return
		}
		imageURL = uploadResult.SecureURL
	}

	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()
	post.ImageURL = imageURL

	_, err = collection.InsertOne(ctx, post)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func UpdatePostFromForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	if title == "" || content == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Error: unable to get userID", http.StatusUnauthorized)
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var existingPost models.Post
	err = collection.FindOne(ctx, bson.M{"_id": postID, "user_id": userID}).Decode(&existingPost)
	if err != nil {
		http.Error(w, "Post not found or does not belong to you", http.StatusForbidden)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":   title,
			"content": content,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": postID}, update)
	if err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/myposts", http.StatusSeeOther)
}

func EditPostPage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Error: unable to get userID", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var post models.Post
	err = collection.FindOne(ctx, bson.M{"_id": postID, "user_id": userID}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found or does not belong to you", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/edit_post.html"))
	if err := tmpl.Execute(w, post); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok || userIDStr == "" {
		http.Error(w, "Error: unable to get userID", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		http.Error(w, "Invalid userID format", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection("posts")
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var existingPost models.Post
	err = collection.FindOne(ctx, bson.M{"_id": postID, "user_id": userID}).Decode(&existingPost)
	if err != nil {
		http.Error(w, "Post not found or does not belong to you", http.StatusForbidden)
		return
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": postID})
	if err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/myposts", http.StatusSeeOther)
}

func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create_post.html"))
	tmpl.Execute(w, nil)
}
