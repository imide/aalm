package player

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/config"
	"github.com/imide/aalm/util/db"
	"log"
	"time"
)

var Drop = cmdutil.Commands{
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

	unknownErrEmbed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An unknown error occurred.", 0xffcc4d)

	// Get all data
	// just gonna import from recruit.go fuck this
	_, teamData, err := GetRecruiterData(s, i)
	if err != nil {
		switch {
		case false:
			return
		default:
			embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An unknown error occurred.", 0xffcc4d)
			cmdutil.SendInteractionResponse(s, i, embed)
			log.Printf("an error occurred during getting data: %s", err)
			return
		}
	}

	playerData, err := db.GetPlayerData(playerID)
	if err != nil {
		log.Printf("Error getting player data: %s", err)
		return
	}

	for _, playerInfo := range teamData.Players {
		if playerInfo.ID == playerID {
			err = db.DropPlayer(&playerData, &teamData)
			if err != nil {
				embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An unknown error occured while editing player/team data.", 0xffcc4d)
				cmdutil.SendInteractionResponse(s, i, embed)
				return
			}

			embed := cmdutil.CreateEmbed("✅ | **Success**", "Player dropped successfully.", 0x00ff00)
			cmdutil.SendInteractionResponse(s, i, embed)

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
				cmdutil.SendInteractionResponse(s, i, unknownErrEmbed)
				return
			}

			_, err = s.ChannelMessageSendEmbed(playerID, dmPlayerEmbed)
			if err != nil {
				log.Printf("Error sending DM to player: %s", err)
				cmdutil.SendInteractionResponse(s, i, unknownErrEmbed)
				return
			}

		}
	}

}
