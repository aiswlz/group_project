package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
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

func EventsPage(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("events")
	if collection == nil {
		http.Error(w, `{"message": "Database connection error"}`, http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sortOrder := -1
	if r.URL.Query().Get("sort") == "old" {
		sortOrder = 1 // сортировка от старых к новым (по возрастанию даты)
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "date", Value: sortOrder}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		http.Error(w, "Error fetching events", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err = cursor.All(ctx, &events); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/events.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, events); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	collection := database.GetCollection("events")
	if collection == nil {
		http.Error(w, `{"message": "Database connection error"}`, http.StatusInternalServerError)
		log.Println("Error: failed to get 'events' collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !middleware.IsAdmin(r) {
		http.Error(w, `{"message": "Access denied"}`, http.StatusForbidden)
		log.Println("Error: user is not an admin")
		return
	}

	var event models.Event
	var imageURL string

	log.Println("Processing request to create event")

	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, `{"message": "Invalid data"}`, http.StatusBadRequest)
			log.Println("JSON decoding error:", err)
			return
		}
		log.Printf("Received event data: %+v", event)
	} else {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, `{"message": "Form processing error"}`, http.StatusBadRequest)
			log.Println("Form parsing error:", err)
			return
		}

		event.Title = r.FormValue("title")
		event.Description = r.FormValue("description")
		log.Println("Received form data - title:", event.Title, "description:", event.Description)

		dateStr := r.FormValue("date")
		eventDate, err := time.Parse("2006-01-02T15:04", dateStr)
		if err != nil {
			http.Error(w, `{"message": "Invalid date format. Use YYYY-MM-DDTHH:MM format"}`, http.StatusBadRequest)
			log.Println("Date parsing error:", err)
			return
		}
		event.Date = eventDate

		file, _, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			uploadResult, err := config.Cloudinary.Upload.Upload(ctx, file, uploader.UploadParams{})
			if err != nil {
				http.Error(w, `{"message": "Image upload error"}`, http.StatusInternalServerError)
				log.Println("Image upload error:", err)
				return
			}
			imageURL = uploadResult.SecureURL
		} else {
			log.Println("Image upload error:", err)
		}
	}

	event.ID = primitive.NewObjectID()
	event.CreatedAt = time.Now()
	event.ImageURL = imageURL

	log.Printf("Creating event: %+v", event)

	if event.Title == "" || event.Description == "" || event.Date.IsZero() {
		http.Error(w, `{"message": "All fields must be filled"}`, http.StatusBadRequest)
		log.Println("Error: one or more required fields are empty.")
		return
	}

	_, err := collection.InsertOne(ctx, event)
	if err != nil {
		http.Error(w, `{"message": "Event creation error"}`, http.StatusInternalServerError)
		log.Println("Error inserting event into database:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Event successfully created!",
		"redirect": "/events",
	})
}

func CreateEventPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/create_event.html")
}
