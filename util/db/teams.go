package db

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
	"time"
)

type TeamInfoOperation int
type TeamDataOperation int

const (
	ChangeLogo TeamInfoOperation = iota
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

var defaultLogo = "" //TODO: add this later
var mu sync.RWMutex
var TeamOptions []*discordgo.ApplicationCommandOptionChoice

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
				Value: t.RoleID,
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

func UpdateTeamInfo(id string, info string, op TeamInfoOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	var update bson.M
	switch op {
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

	err = UpdateTeamOptions()
	if err != nil {
		log.Println(err)
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

	err = UpdateTeamOptions()
	if err != nil {
		log.Println(err)
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

func UpdateMultipleTeamInfo(id string, info bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := client.Database("aafl").Collection("teams") // replace with your actual database name

	for field := range info {
		switch field {
		case "name", "logo", "teamOwner", "coach", "players", "roleId", "playerMax", "wins", "losses", "starsRecruited", "maxStars":
		default:
			return errors.New("invalid field: " + field)
		}
	}

	update := bson.M{"$set": info}

	_, err := c.UpdateOne(ctx, bson.M{"RoleID": id}, update)
	if err != nil {
		return err
	}

	err = UpdateTeamOptions()
	if err != nil {
		log.Println(err)
	}
	return nil
}
