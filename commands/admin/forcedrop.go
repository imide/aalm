package admin

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/db"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var forceDrop = commands.Commands{
	Name:        "forcedrop",
	Description: "Force drop a player from a team",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "player",
			Description: "The player to force drop",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team",
			Description: "The team to force drop the player from",
			Required:    true,
			Choices:     db.TeamOptions,
		},
	},
	Handler: forceDropHandler,
}

func forceDropHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Messages and button stuff:

	var confirmation = &discordgo.MessageEmbed{
		Title: "⚠️ | **Confirmation**",
		Color: 0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Administation",
			//todo: add logo soon
		},
		Timestamp: time.Now().String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Initiated by %s", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
	}

	var cancelEmbed = discordgo.MessageEmbed{
		Title:       "️❌| **Cancelled**",
		Description: "The interaction has been cancelled.",
		Color:       0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Administration",
			//todo: logo
		},
		Timestamp: time.Now().String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Initiated by %s", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
	}

	var successEmbed = discordgo.MessageEmbed{
		Title: "✅ | **Success**",
		Color: 0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Administration",
			//todo: here too
		},
		Timestamp: time.Now().String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Initiated by %s", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
	}

	var forceDropDmEmbed = discordgo.MessageEmbed{
		Title: "⚠️ | **You have been force dropped**",
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Administration",
		},
		Color:     0xffcc4d,
		Timestamp: time.Now().String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "If you believe this is a mistake, please contact the league management.",
		},
	}

	var acceptButton = commands.CreateButton("Accept", discordgo.SuccessButton, "✅", "accept")
	var cancelButton = commands.CreateButton("Cancel", discordgo.DangerButton, "❌", "cancel")
	var confirmRow = discordgo.ActionsRow{Components: []discordgo.MessageComponent{*acceptButton, *cancelButton}}

	// Permission check
	perms, err := s.State.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving your permissions.", 0xffcc4d))
		return
	}

	if perms&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "You do not have permission to use this command.", 0xffcc4d))
		return
	}

	// Retrieve the team data
	teamData, err := db.GetTeamData(i.ApplicationCommandData().Options[1].StringValue())
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the team data.", 0xffcc4d))
		return
	}

	// Retrieve the player data
	playerData, err := db.GetSpecificPlayerData(i.ApplicationCommandData().Options[0].UserValue(s).ID, nil)
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the player data.", 0xffcc4d))
		return
	}

	// Check if the player is already on the team
	if playerData.TeamPlaying != teamData.ID {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "This player is not on the team you specified.", 0xffcc4d))
		return
	}

	// Confirmation
	confirmation.Description = fmt.Sprintf("**Note:** You are doing a potentially destructive action. The following will change by force: \n **Dropping player:** <@%s> from the **%s**", i.ApplicationCommandData().Options[0].UserValue(s).ID, teamData.Name)

	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{confirmation},
			Components: []discordgo.MessageComponent{confirmRow},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	}

	// Send the confirmation
	err = s.InteractionRespond(i.Interaction, &response)
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while sending the confirmation message.", 0xffcc4d))
		return
	}

	// Wait for the user to confirm the action
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			switch i.ApplicationCommandData().ID {
			case "accept":
				// Edit the player data

				playerUpdate := bson.M{
					"team_id":    "",
					"contracted": false,
				}

				err = db.UpdateMultiplePlayerInfo(i.ApplicationCommandData().Options[0].UserValue(s).ID, playerUpdate)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while updating the player data.", 0xffcc4d))
					return
				}

				// Responces
				successEmbed.Description = fmt.Sprintf("The player <@%s> has been force dropped from the team **%s**.", i.ApplicationCommandData().Options[0].UserValue(s).ID, teamData.Name)

				response := discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{&successEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				}

				err = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to the interaction.", 0xffcc4d))
					return
				}

				// DM the user
				dmChannel, err := s.UserChannelCreate(i.ApplicationCommandData().Options[0].UserValue(s).ID)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while creating the DM channel.", 0xffcc4d))
					return
				}

				forceDropDmEmbed.Description = fmt.Sprintf("You have been force dropped from the team **%s**. \n Please note: an **admin** conducted this action, not your coach.", teamData.Name)
				response = discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{&forceDropDmEmbed},
					},
				}

				_, err = s.ChannelMessageSendEmbed(dmChannel.ID, &forceDropDmEmbed)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while sending the DM.", 0xffcc4d))
					return
				}

			case "cancel":
				// Cancel the interaction
				response := discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{&cancelEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				}

				err = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to the interaction.", 0xffcc4d))
					return
				}
			}

		}

	})
}
