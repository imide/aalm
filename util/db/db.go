package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
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
func UpdateMiles(id string, miles int64, op Operation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("test").Collection("miles")

	var update bson.M
	switch op {
	case Replace:
		update = bson.M{"$set": bson.M{"miles": miles}}
	case Add:
		update = bson.M{"$inc": bson.M{"miles": miles}}
	case Subtract:
		update = bson.M{"$inc": bson.M{"miles": -miles}}
	default:
		return errors.New("invalid operation")
	}

	_, err := c.UpdateOne(ctx, bson.M{"userId": id}, update)
	if err != nil {
		return err
	}

	return nil
}

// Exists checks if a user exists in the MongoDB database.
func Exists(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("test").Collection("miles")
	err := c.FindOne(ctx, bson.M{"userId": id}).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// AddUser adds a user to the MongoDB database.
func AddUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("test").Collection("miles")
	_, err := c.InsertOne(ctx, bson.M{"userId": id, "miles": 0})
	if err != nil {
		return err
	}

	return nil
}

// GetMiles retrieves a user's miles from the MongoDB database.
func GetMiles(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("test").Collection("miles")
	var result bson.M
	err := c.FindOne(ctx, bson.M{"userId": id}).Decode(&result)
	if err != nil {
		return 0, err
	}

	miles := result["miles"].(int64)
	return miles, nil
}
