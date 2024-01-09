// Package commands File: commands/utils.go
package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func CreateEmbed(title string, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}

func SendInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
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

func CreateButton(label string, style discordgo.ButtonStyle, emoji string, customId string) *discordgo.Button {
	return &discordgo.Button{
		Label:    label,
		Style:    style,
		CustomID: customId,
		Emoji: discordgo.ComponentEmoji{
			Name: emoji,
		},
	}
}
