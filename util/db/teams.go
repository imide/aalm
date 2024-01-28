package db

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

var (
	defaultLogo = "" // TODO: add this late
	mu          sync.RWMutex
	TeamOptions []*discordgo.ApplicationCommandOptionChoice
)

func UpdateTeamOptions() error {
	teams, err := GetAllTeams()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var options []*discordgo.ApplicationCommandOptionChoice

	for _, team := range teams {
		wg.Add(1)
		go func(t Team) {
			defer wg.Done()
			options = append(options, &discordgo.ApplicationCommandOptionChoice{
				Name:  t.Name,
				Value: t.ID,
			})
		}(team)
	}

	wg.Wait()

	mu.Lock()
	TeamOptions = options
	mu.Unlock()

	return nil
}

func GetAllTeams() ([]Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	cursor, err := c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var teams []Team
	err = cursor.All(ctx, &teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func SaveTeamData(teamData *Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	_, err := c.InsertOne(ctx, teamData)
	if err != nil {
		return err
	}

	err = UpdateTeamOptions()
	if err != nil {
		return err
	}

	return nil
}

func GetTeamData(id string) (Team, error) {
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

func HasManagePermission(userID string, team Team) bool {
	// Check if the user is the team owner
	if userID == team.Owner {
		return true
	}

	// Check if the user is one of the coaches
	for _, coach := range team.Coaches {
		if userID == coach {
			return true
		}
	}

	// If the user is neither the owner nor a coach, they do not have permission
	return false
}
