package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/group-project/aitu-fan/config"
	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/middleware"
	"github.com/group-project/aitu-fan/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GalleryPageData struct {
	GalleryItems []models.GalleryItem
	CurrentSort  string
}

func GetGallery(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("gallery")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sortParam := r.URL.Query().Get("sort")
	var sort bson.D

	switch sortParam {
	case "date_asc":
		sort = bson.D{{"uploaded_at", 1}}
	case "date_desc":
		sort = bson.D{{"uploaded_at", -1}}
	default:
		sort = bson.D{{"uploaded_at", -1}}
		sortParam = "date_desc"
	}

	findOptions := options.Find()
	findOptions.SetSort(sort)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		http.Error(w, "Error fetching gallery", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var gallery []models.GalleryItem
	if err = cursor.All(ctx, &gallery); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	pageData := GalleryPageData{
		GalleryItems: gallery,
		CurrentSort:  sortParam,
	}

	tmpl, err := template.ParseFiles("templates/gallery.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, pageData)
}

func UploadGalleryItem(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("gallery")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID, err := middleware.GetUserID(r)
	if err != nil {
		http.Error(w, "User not authorized", http.StatusUnauthorized)
		return
	}

	r.ParseMultipartForm(10 << 20)
	title := r.FormValue("title")

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "File not uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadResult, err := config.Cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "aitu-fan-gallery",
	})
	if err != nil {
		http.Error(w, "Error uploading image to Cloudinary", http.StatusInternalServerError)
		return
	}

	galleryItem := models.GalleryItem{
		ID:         primitive.NewObjectID(),
		Title:      title,
		URL:        uploadResult.SecureURL,
		UploadedBy: userID,
		UploadedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, galleryItem)
	if err != nil {
		http.Error(w, "Error saving image to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(galleryItem)
}

func GalleryPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/gallery.html"))
	tmpl.Execute(w, nil)
}
