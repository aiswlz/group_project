package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/handlers"
	"github.com/group-project/aitu-fan/models"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var mockUser = models.User{
	Username:     "testuser",
	PasswordHash: "password123",
	Role:         "user",
	CreatedAt:    time.Now(),
}

func TestMain(m *testing.M) {
	err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	m.Run()
}

func TestRegisterUser(t *testing.T) {
	userData := map[string]string{
		"username":      mockUser.Username,
		"password_hash": mockUser.PasswordHash,
	}
	jsonData, err := json.Marshal(userData)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	expectedResponse := `{"message":"Registration successful!"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
}

func TestLoginUser(t *testing.T) {
	loginData := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Error marshalling data: %v", err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	assert.Contains(t, response, "message")
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user_id")
	assert.Contains(t, response, "user_role")

	assert.Equal(t, "Successful!", response["message"])
	assert.Equal(t, "user", response["user_role"])
}

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/home.html")
	})
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestDashboard(t *testing.T) {
	req, err := http.NewRequest("GET", "/dashboard", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Dashboard)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestCreatePostPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/create_post", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreatePostPage)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestCreateEventPage(t *testing.T) {
	req, err := http.NewRequest("GET", "/create_event", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateEventPage)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, status)
	}
}
