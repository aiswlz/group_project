package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"unicode"

	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/models"
	"github.com/group-project/aitu-fan/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendErrorResponse(w, "Error", http.StatusBadRequest)
		return
	}

	if len(user.Username) < 3 || len(user.Username) > 15 {
		sendErrorResponse(w, "Username has to consist of 3-15 symbols", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 6 || len(user.Password) > 8 || !hasSpecialChar(user.Password) {
		sendErrorResponse(w, "Password has to contain 6-8 symbols and at least 1 special symbol", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, "Password hashing error", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     user.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		Role:         "user",
	}

	collection := database.Client.Database("aitu-fan").Collection("users")
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		sendErrorResponse(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successful!"})
}

func hasSpecialChar(s string) bool {
	for _, char := range s {
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			return true
		}
	}
	return false
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	collection := database.Client.Database("aitu-fan").Collection("users")

	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		http.Error(w, "Token generation error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Successful!",
		"token":     tokenString,
		"user_id":   user.ID.Hex(),
		"user_role": user.Role,
	})
}
