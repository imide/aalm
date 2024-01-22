package cmdutil

import "github.com/bwmarrin/discordgo"

type Commands struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Permissions int
}
