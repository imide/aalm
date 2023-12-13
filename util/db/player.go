package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
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
	Position
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

	_, err := c.UpdateOne(ctx, bson.M{"discordId": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePlayerInfo(id string, info string, op PlayerInfoOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("your_database_name").Collection("players") // replace with your actual database name

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
	case Position:
		update = bson.M{"$set": bson.M{"position": info}}
	case ChangeTeam:
		update = bson.M{"$set": bson.M{"teamPlaying": info}}
	default:
		return errors.New("invalid operation")
	}

	_, err := c.UpdateOne(ctx, bson.M{"discordId": id}, update)
	if err != nil {
		return err
	}

	return nil
}
