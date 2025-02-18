package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Date        time.Time          `bson:"date" json:"date"`
	ImageURL    string             `bson:"image_url,omitempty" json:"image_url"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
