package rings

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/db"
	"log"
)

// This command is used to check what rings are available to be awarded to players.

var rings = commands.Commands{
	Name:        "rings",
	Description: "View the rings available to be awarded to players.",
	Handler:     ringsHandler,
}

func ringsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rings, err := db.GetRingsData()
	if err != nil {
		log.Println("Error retrieving rings data,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the rings data.", 0xffcc4d))
		return
	}

	// Embeds

	embed := &discordgo.MessageEmbed{
		Title:       "Rings",
		Description: "Rings available to be awarded to players. The list of available rings is listed below, along with it's data.",
		Color:       0x00ff00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Rings",
			IconURL: "https://cdn.discordapp.com/attachments/881495661745678858/881495690534633728/aafl.png",
		},
	}

	for _, ring := range rings {
		players, err := db.GetPlayersWithRing(ring.ID)
		if err != nil {
			log.Println("Error retrieving players with ring,", err)
			commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the players with ring.", 0xffcc4d))
			return
		}

		playerMentions := ""
		for _, player := range players {
			playerMentions += fmt.Sprintf("<@!%s> ", player.ID)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   ring.Name,
			Value:  fmt.Sprintf("**Description:** %s\n**Players:** %s", ring.Desc, playerMentions),
			Inline: true,
		})

		commands.SendInteractionResponse(s, i, embed)
	}
}
