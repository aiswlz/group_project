package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `json:"username"`
	PasswordHash string             `json:"password_hash"`
	Role         string             `json:"role"`
	CreatedAt    time.Time          `json:"created_at"`
}
