package team

import (
	"fmt"
	"strings"

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
	teamData, err := db.GetTeamData(i.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", fmt.Sprintf("An error occurred while retrieving the team data: %v", err), 0xffcc4d))
		return
	}

	teamRole, err := s.State.Role(i.GuildID, teamData.RoleID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", fmt.Sprintf("An error occurred while retrieving the team role: %v", err), 0xffcc4d))
		return
	}

	var totalStars float32
	fields := make([]*discordgo.MessageEmbedField, 0, len(teamData.Players)+3)
	for _, playerInfo := range teamData.Players {
		totalStars += playerInfo.Stars
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("<@!%s>", playerInfo.ID),
			Value:  fmt.Sprintf("Stars: %f", playerInfo.Stars),
			Inline: true,
		})
	}

	coaches := strings.Join(teamData.Coaches, ", ")
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Coaches",
		Value:  coaches,
		Inline: true,
	})

	ownerData, err := db.GetSpecificPlayerData(teamData.Owner, bson.M{"stars": 1, "_id": 1})
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", fmt.Sprintf("An error occurred while retrieving the player data: %v", err), 0xffcc4d))
		return
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "Franchise Owner",
		Value: fmt.Sprintf("<@!%s>\nStars: %f", ownerData.ID, ownerData.Stars),
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "Stars",
		Value: fmt.Sprintf("Star Cap: %f, Total Stars: %f", teamData.MaxStars, totalStars),
	})

	embed := &discordgo.MessageEmbed{
		Title:       "Team Roster",
		Description: fmt.Sprintf("Roster for %s", teamData.Name),
		Color:       teamRole.Color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: teamData.Logo,
		},
		Fields: fields,
	}

	cmdutil.SendInteractionResponse(s, i, embed)
}
