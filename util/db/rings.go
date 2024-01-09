package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

// Player-Ring functions and cousins

// GetPlayerRings returns the rings of a player.
func GetPlayerRings(id string) ([]PlayerRingData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	var player Player
	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if err != nil {
		return nil, err
	}

	return player.Rings, nil
}

// GetPlayerRingsData returns the data of the rings of a player.
func GetPlayerRingsData(id string) ([]Ring, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("rings")

	playerRings, err := GetPlayerRings(id)
	if err != nil {
		return nil, err
	}

	var ringsData []Ring
	for _, ringID := range playerRings {
		var ring Ring
		err := c.FindOne(ctx, bson.M{"_id": ringID}).Decode(&ring)
		if err != nil {
			return nil, err
		}
		ringsData = append(ringsData, ring)
	}

	return ringsData, nil
}

// GetPlayersWithRing returns the players with a specific ring.
func GetPlayersWithRing(ringID string) ([]Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players")

	cursor, err := c.Find(ctx, bson.M{"rings": bson.M{"$elemMatch": bson.M{"ring_id": ringID}}})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("Error closing cursor,", err)
		}
	}(cursor, ctx)

	var players []Player
	for cursor.Next(ctx) {
		var player Player
		err := cursor.Decode(&player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

// Ring functions and cousins

// CreateRing creates a ring. You must use the Ring struct to create the ring, unlike the other create functions as this one is small and simple.
func CreateRing(ringData Ring) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("rings")

	_, err := c.InsertOne(ctx, ringData)
	if err != nil {
		return err
	}

	return nil
}

// EditRing edits the data of a ring. You must use the Ring struct to edit the data, unlike the other edit functions as this one is small and simple.
func EditRing(id string, ringData Ring) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("rings")

	_, err := c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": ringData})
	if err != nil {
		return err
	}
	return nil
}

// GetRingData returns the data of a single ring. To get the data of multiple rings, use GetRingsData.
func GetRingData(id string) (Ring, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("rings")

	var ring Ring
	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(&ring)
	if err != nil {
		return Ring{}, err
	}

	return ring, nil
}

// GetRingsData returns the data of all rings.
func GetRingsData() ([]Ring, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("rings")

	cursor, err := c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("Error closing cursor,", err)
		}
	}(cursor, ctx)

	var rings []Ring
	for cursor.Next(ctx) {
		var ring Ring
		err := cursor.Decode(&ring)
		if err != nil {
			return nil, err
		}
		rings = append(rings, ring)
	}

	return rings, nil
}
