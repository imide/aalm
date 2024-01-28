package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetPlayerData(id string) (Player, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("players") // replace with your actual database name

	var player Player
	err := c.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
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

func RecruitPlayer(playerData *Player, teamData *Team) error {
	// player
	playerData.TeamPlaying = teamData.ID
	playerData.Contracted = true

	// team
	playerInfo := PlayerInfo{
		ID:    playerData.ID,
		Stars: playerData.Stars,
	}

	teamData.Players = append(teamData.Players, playerInfo)

	// save
	err := SavePlayerData(playerData)
	if err != nil {
		return err
	}

	err = SaveTeamData(teamData)
	if err != nil {
		return err
	}

	return nil
}

func DropPlayer(playerData *Player, teamData *Team) error {
	// player
	playerData.TeamPlaying = ""
	playerData.Contracted = false

	// team
	index := -1
	for i, playerInfo := range teamData.Players {
		if playerInfo.ID == playerData.ID {
			index = i
			break
		}
	}

	if index != -1 {
		teamData.Players = append(teamData.Players[:index], teamData.Players[index+1:]...)
	}

	// save
	err := SavePlayerData(playerData)
	if err != nil {
		return err
	}

	err = SaveTeamData(teamData)
	if err != nil {
		return err
	}

	return nil
}
