package player

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/db"
	"log"
)

var ringcheck = commands.Commands{
	Name:        "ringcheck",
	Description: "gets the amount of rings of a player so you can make fun of them. or to flex.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "player",
			Description: "The player to view the rings of",
			Required:    true,
		},
	},
	Handler: ringcheckHandler,
}

func ringcheckHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	player, err := s.User(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	if err != nil {
		log.Println("Error retrieving user,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the user.", 0xffcc4d))
		return
	}

	playerRings, err := db.GetPlayerRingsData(player.ID)
	if err != nil {
		log.Println("Error retrieving player rings,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the player rings.", 0xffcc4d))
		return
	}

	switch len(playerRings) {
	case 0:
		commands.SendInteractionResponse(s, i, &discordgo.MessageEmbed{
			Title:       "lmao.",
			Description: fmt.Sprintf("the absolute loser <@!%s> has no rings. make fun of this user.", player.ID),
			Color:       0x800000,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("AAFL Ring Check"),
				IconURL: player.AvatarURL(""),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://tenor.com/view/lachen-gif-5717730327661866096",
			},
		})
	default:
		commands.SendInteractionResponse(s, i, &discordgo.MessageEmbed{
			Title:       "damn ok.",
			Description: fmt.Sprintf("<@!%s> has %s rings. seems pretty good. \n **To view more in-depth information, please run the /player command.**", player.ID, len(playerRings)),
			Color:       0x008000,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintf("AAFL Ring Check"),
				IconURL: player.AvatarURL(""),
			},
		})
	}
}
