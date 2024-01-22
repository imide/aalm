package db

import (
	"context"
	"github.com/imide/aalm/util/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

func Init() {
	var err error

	cfg := config.Cfg

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
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
