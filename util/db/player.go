package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type PlayerStarOperation int
type PlayerInfoOperation int

const (
	Set PlayerStarOperation = iota
	Add
	Remove
)

const (
	Suspend PlayerInfoOperation = iota
	Unsuspend
	Contract
	Drop
	AddRing
	RemoveRing
	ChangeTeam
)

func UpdatePlayerStars(id string, stars int, op PlayerStarOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("your_database_name").Collection("players") // replace with your actual database name

	var update bson.M
	switch op {
	case Set:
		update = bson.M{"$set": bson.M{"stars": stars}}
	case Add:
		update = bson.M{"$inc": bson.M{"stars": stars}}
	case Remove:
		update = bson.M{"$inc": bson.M{"stars": -stars}}
	default:
		return errors.New("invalid operation")
	}

	_, err := c.UpdateOne(ctx, bson.M{"ID": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePlayerInfo(id string, info string, op PlayerInfoOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	var update bson.M
	switch op {
	case Suspend:
		update = bson.M{"$set": bson.M{"isSuspended": true, "suspensionExpires": info}}
	case Unsuspend:
		update = bson.M{"$set": bson.M{"isSuspended": false}}
	case Contract:
		update = bson.M{"$set": bson.M{"contracted": true}}
	case Drop:
		update = bson.M{"$set": bson.M{"contracted": false}}
	case ChangeTeam:
		update = bson.M{"$set": bson.M{"teamPlaying": info}}
	case AddRing:
		update = bson.M{"$push": bson.M{"rings": info}}
	case RemoveRing:
		update = bson.M{"$pull": bson.M{"rings": info}}

	default:
		return errors.New("invalid operation")
	}

	_, err := c.UpdateOne(ctx, bson.M{"ID": id}, update)
	if err != nil {
		return err
	}

	return nil
}

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

func UpdateMultiplePlayerInfo(id string, info bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	for field := range info {
		switch field {
		case "teamPlaying", "rings", "contracted", "seasonsContracted", "isSuspended", "suspensionExpires", "stars":

		default:
			return errors.New("invalid field: " + field)
		}
	}

	update := bson.M{"$set": info}

	_, err := c.UpdateOne(ctx, bson.M{"ID": id}, update)
	if err != nil {
		return err
	}

	return nil
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
