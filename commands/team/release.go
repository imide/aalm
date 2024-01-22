package team

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/db"
)

var Release = cmdutil.Commands{
	Name:        "release",
	Description: "Releases a player from your team.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "player",
			Description: "The player to release.",
			Required:    true,
		},
	},
	Handler: releaseHandler,
}

func releaseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Get data
	_, coachTeam, err := getCoachData(s, i)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving your team data.", 0xffcc4d))
		return
	}

	// Get player data
	playerData, err := db.GetPlayerData(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the player data.", 0xffcc4d))
		return
	}

	// Check if player is on team

	if playerData.TeamPlaying != coachTeam.ID {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "This player is not on your team.", 0xffcc4d))
		return
	}

	// Important variables

	var confirmation = &discordgo.MessageEmbed{
		Title:       "⚠️ | **Warning**",
		Description: fmt.Sprintf("Are you sure you want to release <@!%s> from your team?", playerData.ID),
		Color:       0xffcc4d,
	}

	var success = &discordgo.MessageEmbed{
		Title:       "✅ | **Success**",
		Description: fmt.Sprintf("<@!%s> has been released from your team. You may view your team's new roster via /roster.", playerData.ID),
		Color:       0x00ff00,
	}

	var playerDm = &discordgo.MessageEmbed{
		Title:       "⚠️ | **Warning**",
		Description: fmt.Sprintf("You have been released from %s.", coachTeam.Name),
		Color:       0xffcc4d,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: coachTeam.Logo,
		},
	}

	var acceptButton = cmdutil.CreateButton("Accept", discordgo.SuccessButton, "✅", "accept")

	var denyButton = cmdutil.CreateButton("Deny", discordgo.DangerButton, "❌", "deny")

	var confirmRow = discordgo.ActionsRow{Components: []discordgo.MessageComponent{*acceptButton, *denyButton}}

	// Actual logic

	message := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{confirmation},
			Flags:  discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				confirmRow,
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, &message)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while sending the confirmation message.", 0xffcc4d))
		return
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			switch i.ApplicationCommandData().ID {
			case "accept":
				// Remove player from team

				playerData.TeamPlaying = ""
				playerData.Contracted = false

				err = db.SavePlayerData(&playerData)

				if err != nil {
					cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while removing the player from the team.", 0xffcc4d))
					return
				}

				// Send success message
				cmdutil.SendInteractionResponse(s, i, success)

				// Send dm to player
				_, err = s.UserChannelCreate(playerData.ID)
				if err != nil {
					cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while creating the user channel.", 0xffcc4d))
					return
				}

				_, err = s.ChannelMessageSendEmbed(playerData.ID, playerDm)
				if err != nil {
					cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while sending the dm to the player.", 0xffcc4d))
					return
				}
			case "deny":
				cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "The action has been cancelled.", 0xffcc4d))
				return
			}
		}
	})
}

func getCoachData(s *discordgo.Session, i *discordgo.InteractionCreate) (db.Player, db.Team, error) {
	coachData, err := db.GetPlayerData(i.Member.User.ID)
	if err != nil {
		return db.Player{}, db.Team{}, err
	}

	coachTeam, _ := db.GetTeamData(coachData.TeamPlaying)

	if !db.HasManagePermission(i.Member.User.ID, coachTeam) {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "You do not have permission to manage this team. Please try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return db.Player{}, db.Team{}, nil
	}

	return coachData, coachTeam, nil
}
