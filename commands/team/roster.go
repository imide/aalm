package team

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/db"
	"go.mongodb.org/mongo-driver/bson"
)

var Roster = cmdutil.Commands{
	Name:        "roster",
	Description: "Views the roster of the team selected.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team",
			Description: "The team to view the roster of.",
			Required:    true,
			Choices:     db.TeamOptions,
		},
	},
	Handler: rosterHandler,
}

func rosterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Retrieve the team data
	teamData, err := db.GetTeamData(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		log.Println("Error retrieving team data,", err)
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the team data.", 0xffcc4d))
		return
	}
	teamRole, err := s.State.Role(i.GuildID, teamData.RoleID)
	if err != nil {
		log.Println("Error retrieving team role,", err)
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the team role.", 0xffcc4d))
		return
	}

	// Create a new MessageEmbed
	embed := &discordgo.MessageEmbed{
		Title:       "Team Roster",
		Description: fmt.Sprintf("Roster for %s", teamData.Name),
		Color:       teamRole.Color, // Set the color to the role color
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: teamData.Logo, // Set the thumbnail to the team's logo
		},
	}

	// Add a field for each player
	var totalStars float32 = 0
	for _, playerInfo := range teamData.Players {
		totalStars += playerInfo.Stars
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("<@!%s>", playerInfo.ID),
			Value:  fmt.Sprintf("Stars: %f", playerInfo.Stars),
			Inline: true,
		})
	}

	// Add a field for the coaches
	var coaches string
	for _, coachID := range teamData.Coaches {
		coaches += fmt.Sprintf("<@!%s>, ", coachID)
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Coaches",
		Value:  coaches,
		Inline: true,
	})

	// Add a field for the franchise owner
	ownerData, err := db.GetSpecificPlayerData(teamData.Owner, bson.M{"stars": 1, "_id": 1})
	if err != nil {
		log.Println("Error retrieving player data,", err)
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the player data.", 0xffcc4d))
		return
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name: "Franchise Owner",
		Value: fmt.Sprintf(
			"<@!%s>\nStars: %f", ownerData.ID, ownerData.Stars),
	})

	// Add a field for the star cap and the total stars
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "Stars",
		Value: fmt.Sprintf("Star Cap: %f, Total Stars: %f", teamData.MaxStars, totalStars),
	})

	// Send the embed as a response to the interaction
	cmdutil.SendInteractionResponse(s, i, embed)
}
