package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Like struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID primitive.ObjectID `bson:"post_id" json:"post_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}
