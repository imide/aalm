package admin

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/auditlog"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/db"
	"log"
	"math"
	"time"
)

var suspend = &commands.Commands{
	Name:        "suspend",
	Description: "Suspends a user for a given amount of time. Suspensions are stored in the player's document in the database indefinitely, unless manually removed/overridden in database.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to suspend",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "The reason for the suspension",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "duration",
			Description: "The duration of the suspension WILL DEFAULT TO FOREVER IF NOT SPECIFIED. Example: 4h30m is for 4 hours 30 minutes. 7d is a week.",
			Required:    false,
		},
	},
	Handler: suspendHandler,
}

// suspend a player for a given amount of time. Suspensions are stored in the player's document in the database indefinitely, unless manually removed/overridden in database.
func suspendHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	// Messages and button stuff:

	var confirmation = &discordgo.MessageEmbed{
		Title: "⚠️ | **Confirmation**",
		Color: 0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Suspensions",
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
		Description: "The suspension has been cancelled.",
		Color:       0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Suspensions",
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
			Name: "Creatine | Suspensions",
			//todo: here too
		},
		Timestamp: time.Now().String(),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Initiated by %s", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
	}

	var suspendDmEmbed = discordgo.MessageEmbed{
		Title: "⚠️ | **You have been suspended**",
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Creatine | Suspensions",
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

	// Logic:
	// Retrieve the user data
	userData, err := db.GetPlayerData(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the user data.", 0xffcc4d))
		return
	}

	// Check if user is already suspended

	if userData.Suspension != nil && userData.Suspension.Active == true {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "This user is already suspended.", 0xffcc4d))
		return
	}

	// Retrieve the duration
	var duration time.Duration
	if i.ApplicationCommandData().Options[2].StringValue() == "" {
		duration = math.MaxInt64
		confirmation.Description = fmt.Sprintf("Are you sure you want to suspend <@%s> for **forever**?", i.ApplicationCommandData().Options[0].UserValue(s).ID)
	} else {
		duration, err = time.ParseDuration(i.ApplicationCommandData().Options[2].StringValue())
		if err != nil {
			commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while parsing the duration.", 0xffcc4d))
			return
		}
		confirmation.Description = fmt.Sprintf("Are you sure you want to suspend <@%s> for **%s**?", i.ApplicationCommandData().Options[0].UserValue(s).ID, duration.String())
	}

	// Prompt for confirmation
	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{confirmation},
			Components: []discordgo.MessageComponent{confirmRow},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	}

	err = s.InteractionRespond(i.Interaction, &response)
	if err != nil {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to your interaction.", 0xffcc4d))
		log.Println("Error responding to interaction,", err)
		return
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			switch i.ApplicationCommandData().ID {
			case "accept":
				// Prepare the suspension
				var suspension = db.Suspension{
					WhoSuspended: i.Member.User.ID,
					Reason:       i.ApplicationCommandData().Options[1].StringValue(),
					Started:      time.Now(),
					Ends:         time.Now().Add(duration).Unix(),
					Active:       true,
				}

				// Update the player's document
				err = db.UpdatePlayerInfo(i.ApplicationCommandData().Options[0].UserValue(s).ID, suspension, db.SuspendOrUnsuspend)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while updating the player's document.", 0xffcc4d))
					return
				}

				// Send success embed
				response := discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{&successEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				}

				err = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to your interaction.", 0xffcc4d))
					log.Println("Error responding to interaction,", err)
					return
				}

				// Send DM to user

				if duration == math.MaxInt64 {
					suspendDmEmbed.Description = fmt.Sprintf("You have been suspended from the AAFL for **forever**. \n- Reason: **%s**\n**This is a permanent suspension that CANNOT be appealed unless specified.", suspension.Reason)
				} else {
					unixSuspendTime := time.Now().Add(duration).Unix()
					suspendDmEmbed.Description = fmt.Sprintf("You have been suspended from the AAFL for **%s**. \n- Reason: **%s** \n **You will be unsuspended in <t:%s:R> ", duration.String(), suspension.Reason, unixSuspendTime)
				}

				_, err = s.ChannelMessageSendEmbed(i.ApplicationCommandData().Options[0].UserValue(s).ID, &suspendDmEmbed)
				if err != nil {
					log.Println("Error sending message,", err)
					return
				}

				err = auditlog.LogSuspension(s, userData, auditlog.Suspend)
				if err != nil {
					log.Println("Error logging suspension,", err)
					return
				}

				return
			case "cancel":
				// Send cancel embed
				response := discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{&cancelEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				}

				err = s.InteractionRespond(i.Interaction, &response)
				if err != nil {
					commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to your interaction.", 0xffcc4d))
					log.Println("Error responding to interaction,", err)
					return
				}

			}

		}

	})
}
