package misc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
)

var goon = commands.Commands{
	Name:        "goon",
	Description: "Goon.",
	Options:     nil,
	Handler:     goonHandler,
}

func goonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "im literally gooning rn cause its sigma",
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to the interaction.", 0xffcc4d))
		return
	}
}
