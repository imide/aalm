package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type Commands struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CmdMap = map[string]Commands{}

func Register(s *discordgo.Session, guildID string) {
	// Fetch existing commands
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		log.Fatalf("Failed to fetch existing commands: %v", err)
	}

	// Create a map of local commands for easy lookup
	localCommands := make(map[string]Commands)
	for _, cmd := range CmdMap {
		localCommands[cmd.Name] = cmd
	}

	// Iterate over existing commands
	for _, cmd := range existingCommands {
		// If an existing command does not exist in local commands, delete it
		if _, exists := localCommands[cmd.Name]; !exists {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
			if err != nil {
				log.Fatalf("Failed to delete command: %v", err)
			}
		}
	}

	// Register local commands
	for _, command := range CmdMap {
		createCmd(s, guildID, command)
	}
}

func createCmd(s *discordgo.Session, guildID string, cmd Commands) {
	_, err := s.ApplicationCommandCreate(s.State.Application.ID, guildID, &discordgo.ApplicationCommand{
		Name:        cmd.Name,
		Description: cmd.Description,
		Options:     cmd.Options,
	})
	if err != nil {
		log.Println("Err creating cmd,", err)
	}
}
