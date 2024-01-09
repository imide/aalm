package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client

var uri = os.Getenv("MONGODB_URI")

func init() {
	var err error
	clientOptions := options.Client().ApplyURI(uri)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
}

// UpdateMiles updates a user's miles in the MongoDB database based on the provided operation.
