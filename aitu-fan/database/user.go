package database

import (
	"context"
	"log"

	"github.com/group-project/aitu-fan/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection() *mongo.Collection {
	return Client.Database("aitu-fan").Collection("users")
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	collection := getUserCollection()

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error getting users:", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(userID string) (*models.User, error) {
	collection := getUserCollection()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Error converting ID:", err)
		return nil, err
	}

	var user models.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Println("Error finding user:", err)
		return nil, err
	}

	return &user, nil
}

func DeleteUserByID(userID string) error {
	collection := getUserCollection()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("Error converting ID:", err)
		return err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Println("Error deleting user:", err)
	}
	return err
}
