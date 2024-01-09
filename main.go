package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/config"
	"log"
	"os"
	"os/signal"
)

// main initializes the Discord bot by loading environment variables,
// starting the Discord session, and testing the database connection.
func main() {

	// Standard env loading
	err := config.Init()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	// Force change timezone to EST so people dont get mad at me :(
	err = os.Setenv("TZ", "America/New_York")
	if err != nil {
		return
	}

	// Start the Discord session
	startDiscordSession()
}

// startDiscordSession creates and manages a Discord session.
func startDiscordSession() {
	token := fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN"))
	session, err := discordgo.New(token)
	if err != nil {
		log.Fatalf("Creating session error: %v", err)
		return
	}

	if session == nil {
		log.Fatalf("Session is nil")
		return
	}

	if err := session.Open(); err != nil {
		log.Fatalf("Opening connection error: %v", err)
	}

	log.Println("Connection opened")

	go commands.Register(session, os.Getenv("GUILD_ID"))
	go registerCmdHandlers(session)

	log.Println("Bot is running. Press CTRL-C to exit.")
	waitForInterrupt()

	if err := session.Close(); err != nil {
		log.Fatalf("Closing session error: %v", err)
	}
}

// registerCmdHandlers registers command handlers to the session.
func registerCmdHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand && i.Data != nil {
			data, ok := i.Data.(discordgo.ApplicationCommandInteractionData)
			if !ok {
				log.Println("Error casting i.Data to ApplicationCommandInteractionData")
				return
			}
			if command, exists := commands.CmdMap[data.Name]; exists {
				command.Handler(s, i)
			}
		}
	})
}

// waitForInterrupt waits for CTRL-C or other term signal.
func waitForInterrupt() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
