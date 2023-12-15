package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type TeamInfoOperation int
type TeamDataOperation int

const (
	Rename TeamInfoOperation = iota
	ChangeLogo
	Owner
	Coaches
	Players
	Role
)

const (
	MaxPlayers TeamDataOperation = iota
	Wins
	Losses
	Stars
	Starcap
)

var defaultLogo = "" //add this later

func UpdateTeamInfo(id string, info string, op TeamInfoOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	var update bson.M
	switch op {
	case Rename:
		update = bson.M{"$set": bson.M{"name": info}}
	case ChangeLogo:
		update = bson.M{"$set": bson.M{"logo": info}}
	case Owner:
		update = bson.M{"$set": bson.M{"teamOwner": info}}
	case Coaches:
		update = bson.M{"$set": bson.M{"coach": info}}
	case Players:
		update = bson.M{"$set": bson.M{"players": info}}
	case Role:
		update = bson.M{"$set": bson.M{"roleId": info}}
	default:
		return errors.New("invalid operation")
	}

	_, err := c.UpdateOne(ctx, bson.M{"RoleID": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTeamData(id string, data int, op TeamDataOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	var update bson.M
	switch op {
	case MaxPlayers:
		update = bson.M{"$set": bson.M{"playerMax": data}}
	case Wins:
		update = bson.M{"$set": bson.M{"wins": data}}
	case Losses:
		update = bson.M{"$set": bson.M{"losses": data}}
	case Stars:
		update = bson.M{"$set": bson.M{"starsRecruited": data}}
	case Starcap:
		update = bson.M{"$set": bson.M{"maxStars": data}}
	default:
		return errors.New("invalid operation")
	}

	_, err := c.UpdateOne(ctx, bson.M{"RoleID": id}, update)
	if err != nil {
		return err
	}

	return nil
}

func GetTeamInfo(id string) (Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	var team Team
	err := c.FindOne(ctx, bson.M{"RoleID": id}).Decode(&team)
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func AddTeamPlaceholder(team Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	_, err := c.InsertOne(ctx, team)
	if err != nil {
		return err
	}

	return nil
}

func AddTeam(id string, name string, owner string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	team := Team{
		Name:           name,
		Logo:           defaultLogo,
		TeamOwner:      owner,
		Coach:          []string{},
		Players:        []string{},
		PlayerMax:      0,
		RoleID:         id,
		Wins:           0,
		Losses:         0,
		StarsRecruited: 0,
		MaxStars:       30,
	}

	_, err := c.InsertOne(ctx, team)
	if err != nil {
		return err
	}

	playerUpdate := bson.M{
		"teamPlaying": name,
		"position":    "owner",
		// todo: anything else i may be forgettin
	}

	err = UpdateMultiplePlayerInfo(owner, playerUpdate)
	if err != nil {
		return err
	}

	return nil
}
