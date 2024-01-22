package misc

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
)

var Goon = cmdutil.Commands{
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
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to the interaction.", 0xffcc4d))
		return
	}
}
