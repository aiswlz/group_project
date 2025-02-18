package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GalleryItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL        string             `bson:"url" json:"url"`
	Title      string             `bson:"title,omitempty" json:"title"`
	UploadedBy string             `bson:"uploaded_by" json:"uploaded_by"`
	UploadedAt time.Time          `bson:"uploaded_at" json:"uploaded_at"`
}
