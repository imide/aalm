package player

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/config"
	"github.com/imide/aalm/util/db"
	"log"
	"time"
)

var drop = &commands.Commands{
	Name:        "drop",
	Description: "Drops a player from your team.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "player",
			Description: "The player to drop",
			Required:    true,
		},
	},
	Handler: dropHandler,
}

func dropHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// shortcut
	var cfg = config.Cfg
	var playerID = i.ApplicationCommandData().Options[0].UserValue(s).ID

	// reused messages

	unknownErrEmbed := commands.CreateEmbed("⚠️ | **Warning**", "An unknown error occurred.", 0xffcc4d)

	// Get all data
	// just gonna import from recruit.go fuck this
	_, teamData, err := GetRecruiterData(s, i)
	if err != nil {
		switch {
		case false:
			return
		default:
			embed := commands.CreateEmbed("⚠️ | **Warning**", "An unknown error occurred.", 0xffcc4d)
			commands.SendInteractionResponse(s, i, embed)
			log.Printf("an error occurred during getting data: %s", err)
			return
		}
	}

	playerData, err := db.GetPlayerData(playerID)
	if err != nil {
		log.Printf("Error getting player data: %s", err)
		return
	}

	for index, id := range teamData.Players {
		if id == playerID {
			// Remove the player from the list
			teamData.Players = append(teamData.Players[:1], teamData.Players[index+1:]...)

			// Update player data
			playerData.TeamPlaying = ""

			// Save
			err := db.SaveTeamData(&teamData)
			if err != nil {
				log.Printf("Error saving team data: %s", err)
				commands.SendInteractionResponse(s, i, unknownErrEmbed)
				return
			}

			err = db.SavePlayerData(&playerData)
			if err != nil {
				log.Printf("Error saving player data: %s", err)
				embed := commands.CreateEmbed("⚠️ | **Warning**", "An unknown error occurred.", 0xffcc4d)
				commands.SendInteractionResponse(s, i, embed)
				return
			}

			embed := commands.CreateEmbed("✅ | **Success**", "Player dropped successfully.", 0x00ff00)
			commands.SendInteractionResponse(s, i, embed)

			// Send DM to player
			var teamRole, _ = s.State.Role(cfg.GuildID, teamData.RoleID)

			dmPlayerEmbed := &discordgo.MessageEmbed{
				Title:       "⚠️ | **Warning**",
				Timestamp:   time.Now().String(),
				Color:       teamRole.Color,
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: teamData.Logo},
				Author:      &discordgo.MessageEmbedAuthor{Name: "Creatine | Transactions", IconURL: ""},
				Description: fmt.Sprintf("The coaches of your team **(%s)** have dropped you from their roster. \n You are now a free agent and may be recruited by other teams.", teamData.Name),
				Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("If you believe this is a mistake, please contact the coaches of your team. Dropped at: %s", time.Now().String())},
			}

			_, err = s.UserChannelCreate(playerID)
			if err != nil {
				log.Printf("Error creating DM channel: %s", err)
				commands.SendInteractionResponse(s, i, unknownErrEmbed)
				return
			}

			_, err = s.ChannelMessageSendEmbed(playerID, dmPlayerEmbed)
			if err != nil {
				log.Printf("Error sending DM to player: %s", err)
				commands.SendInteractionResponse(s, i, unknownErrEmbed)
				return
			}

			//TODO: audit
		}
	}

}
