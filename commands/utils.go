// Package commands File: commands/utils.go
package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/util/bloxlink"
	"log"
)

// createOptions creates and returns the command options dynamically.
func createOptions(optionTypes []discordgo.ApplicationCommandOptionType, optionNames []string, optionDescriptions []string) []*discordgo.ApplicationCommandOption {
	var options []*discordgo.ApplicationCommandOption

	for i := 0; i < len(optionTypes); i++ {
		option := &discordgo.ApplicationCommandOption{
			Type:        optionTypes[i],
			Name:        optionNames[i],
			Description: optionDescriptions[i],
			Required:    true,
		}
		options = append(options, option)
	}

	return options
}

func getUser(s *discordgo.Session, i *discordgo.InteractionCreate) (string, string) {
	userId := i.ApplicationCommandData().Options[0].UserValue(s).ID
	robloxId, _ := bloxlink.GetRobloxId(userId)
	return userId, robloxId
}

func isBloxlinkFailed(robloxId string) bool {
	return robloxId == "Bloxlink Connection Failed"
}

func handleBloxlinkFailure(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := createEmbed("❌ | **Cancelled**", "User not verified with Bloxlink.", 0xFF0000)
	sendInteractionResponse(s, i, embed)
}

func handleUserDoesNotExist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := createEmbed("⚠️ | **Warning**", "The user specified is not registered in the database. Please create the user with /adduser before continuing.", 0xffcc4d)
	sendInteractionResponse(s, i, embed)
}

func promptUserCreation(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	embed := createEmbed("⚠️ | **Warning**", "Would you like to create a new user?", 0xffcc4d)

	confirmed := sendConfirmation(s, i, embed)

	switch confirmed {
	case true:
		log.Println("interaction user make confirmed, creating user")
		return true

	case false:
		log.Println("interaction user make denied, cancelling")
		return false
	}
	return false
}

func createEmbed(title string, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}

func sendInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println("Error responding to interaction,", err)
	}
}

func confirmHandle(s *discordgo.Session) bool {
	result := make(chan bool)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}
		switch i.MessageComponentData().CustomID {
		case "confirm":
			result <- true
		case "deny":
			result <- false
		}
	})
	return <-result
}

func sendConfirmation(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) bool {
	confirm := createButton("Confirm", discordgo.SuccessButton, "✅", "confirm")
	deny := createButton("Deny", discordgo.DangerButton, "❌", "deny")
	components := discordgo.ActionsRow{Components: []discordgo.MessageComponent{*confirm, *deny}}
	sendInteractionWithComponent(s, i, embed, []discordgo.MessageComponent{components})
	return confirmHandle(s)
}

func sendInteractionWithComponent(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Println("Error responding to interaction,", err)
	}
}

func sendFollowupComponent(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components []discordgo.MessageComponent) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Flags:      discordgo.MessageFlagsEphemeral,
		Components: components,
	})
	if err != nil {
		log.Println("Error responding to interaction,", err)
	}
}

func sendFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Println("Error responding to interaction,", err)
	}

}

func createButton(label string, style discordgo.ButtonStyle, emoji string, customId string) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    style,
		CustomID: customId,
		Emoji: discordgo.ComponentEmoji{
			Name: emoji,
		},
	}
}
