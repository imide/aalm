package trade

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/db"
)

var Trade = cmdutil.Commands{
	Name:        "trade",
	Description: "Start a trade with another team",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "start",
			Description: "Start a trade with another team",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "team",
					Description: "The team you want to trade with",
					Required:    true,
					Choices:     db.TeamOptions,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "accept",
			Description: "Accept a trade with another team",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "team",
					Description: "The team you want to accept the trade with",
					Required:    true,
					Choices:     db.TeamOptions,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "decline",
			Description: "Decline a trade with another team",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "team",
					Description: "The team you want to decline the trade with",
					Required:    true,
					Choices:     db.TeamOptions,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "negociate",
			Description: "Negotiate an existing trade with another team",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "team",
					Description: "The team you want to negotiate the trade with",
					Required:    true,
					Choices:     db.TeamOptions,
				},
			},
		},
	},
}
