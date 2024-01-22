package cmdutil

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/util/db"
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
		Emoji: &discordgo.ComponentEmoji{
			Name: emoji,
		},
	}
}

func UserDoesntExist(uid string) (db.Player, error) {
	// Create a new player with the default values but append the ID

	playerData := db.Player{
		ID:          uid,
		Stars:       0,
		TeamPlaying: "",
		Contracted:  false,
	}

	// Save the player data
	err := db.SavePlayerData(&playerData)
	if err != nil {
		log.Printf("Error saving player data: %s", err)
		return db.Player{}, err
	}

	return playerData, nil
}
