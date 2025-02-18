package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectMongoDB() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("Error connecting to MongoDB: %v", err)
	}

	fmt.Println("Connection to MongoDB is successful!")
	Client = client
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("aitu-fan").Collection(collectionName)
}
