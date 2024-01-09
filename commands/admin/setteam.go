package admin

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/db"
	"log"
)

var setTeam = commands.Commands{
	Name:        "setteam",
	Description: "Set the team of a user",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to set the team of",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team",
			Description: "The team to set the user to",
			Required:    true,
			Choices:     db.TeamOptions,
		},
	},
	Handler: setTeamHandler,
}

func setTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Permission check

	// Retrieve the team data
	teamData, err := db.GetTeamData(i.ApplicationCommandData().Options[1].StringValue())
	if err != nil {
		log.Println("Error retrieving team data,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the team data.", 0xffcc4d))
		return
	}
	teamRole, err := s.State.Role(i.GuildID, teamData.RoleID)
	if err != nil {
		log.Println("Error retrieving team role,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the team role.", 0xffcc4d))
		return
	}

	// Retrieve the user data
	userData, err := db.GetPlayerData(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	if err != nil {
		log.Println("Error retrieving user data,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the user data.", 0xffcc4d))
		return
	}

	// Messages and button stuff:

	var confirmation = &discordgo.MessageEmbed{
		Title:       "⚠️ | **Wait!**",
		Description: fmt.Sprintf("**Note:** You are doing a potentially destructive action. The following will change by force: \n "),
		Color:       0xffcc4d,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  fmt.Sprintf("<@%s> to be force traded to:", i.ApplicationCommandData().Options[0].UserValue(s).ID),
				Inline: true,
			},
			{
				Name:   "Team",
				Value:  fmt.Sprintf("<@&%s> ", i.ApplicationCommandData().Options[1].StringValue()),
				Inline: true,
			},
		},
	}

	var playerAlert = &discordgo.MessageEmbed{
		Title:       "⚠️ | **Warning**",
		Description: fmt.Sprintf("You have been forced traded to the the **%s** by an admin. Your data has automatically been edited. Please join the team (discord)[discord.gg/%i", teamData.Name, teamData.DiscordInvite),
		Color:       teamRole.Color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: teamData.Logo,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("AAFL Recruitment"),
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
	}

	var acceptButton = commands.CreateButton("Accept", discordgo.SuccessButton, "✅", "accept")

	var denyButton = commands.CreateButton("Deny", discordgo.DangerButton, "❌", "deny")

	var confirmRow = discordgo.ActionsRow{Components: []discordgo.MessageComponent{*acceptButton, *denyButton}}

	// Standard checks for my sanity

	if userData.Contracted == false {
		log.Println("Error retrieving user data,", err)
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "This user is not contracted. The user is most likely a free agent, which will not work with this command. Please try again later.", 0xffcc4d))
		return
	}

	if userData.TeamPlaying == i.ApplicationCommandData().Options[1].StringValue() {
		commands.SendInteractionResponse(s, i, commands.CreateEmbed("⚠️ | **Warning**", "This user is already on this team.", 0xffcc4d))
		return
	}

	// Asks the user to confirm the action

	message := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{confirmation},
			Flags:  discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				confirmRow,
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, message)
	if err != nil {
		log.Println("Error responding to interaction,", err)
		return
	}

	// Waits for the user to confirm the action
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			switch i.ApplicationCommandData().ID {
			case "accept":
				//TODO: audit log
				db.UpdatePlayerInfo(i.ApplicationCommandData().Options[0].UserValue(s).ID, teamData.Name, db.ChangeTeam)
				dmChannel, err := s.UserChannelCreate(i.ApplicationCommandData().Options[0].UserValue(s).ID)
				if err != nil {
					log.Println("Error creating DM channel,", err)
					return
				}
				_, err = s.ChannelMessageSendEmbed(dmChannel.ID, playerAlert)
				if err != nil {
					log.Println("Error sending message,", err)
					return
				}
				commands.SendInteractionResponse(s, i, commands.CreateEmbed("✅ | **Success**", "The user has been force traded.", 0x00ff00))
				return

			}
		}
	})
}
