package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	ImageURL  string             `bson:"image_url,omitempty" json:"image_url"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username,omitempty" json:"username"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
