package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func GetPlayerData(id string) (Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	var player Player
	err := c.FindOne(ctx, bson.M{"ID": id}).Decode(&player)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func GetSpecificPlayerData(id string, projection bson.M) (Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	var player Player
	err := c.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(projection)).Decode(&player)
	if err != nil {
		return Player{}, err
	}

	return player, nil
}

func SavePlayerData(playerData *Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	_, err := c.UpdateOne(ctx, bson.M{"_id": playerData.ID}, bson.M{"$set": playerData})
	if err != nil {
		return err
	}

	return nil
}
