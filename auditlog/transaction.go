package auditlog

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/util/config"
	"github.com/imide/aalm/util/db"
	"log"
	"time"
)

// This will contain the audit log functions relating to all transactions of players (trades, cuts, suspensions, etc.)

type TransactionType int

const (
	Trade TransactionType = iota
	Cut
	Contract
	Drop
)

//TODO: add input params cause ik i need more im just doing contracting for now

func LogTransaction(s *discordgo.Session, action TransactionType, player db.Player, team db.Team) {
	// Variable shortcuts cause im fucking lazy and i dont care
	var cfg = config.Cfg
	var playerID, playerStars = player.ID, player.Stars // player shortcuts (add when needed)
	var teamRID, teamLName, teamPCap, teamSCap, teamLogo, teamSName, teamEID = team.RoleID, team.Name, team.PlayerMax, team.MaxStars, team.Logo, team.ID, team.EmojiID
	var teamRole, _ = s.State.Role(cfg.GuildID, teamRID)

	// Star and player cap
	var teamPlayersStars float32
	var teamPlayersNum = len(team.Players)

	for _, player := range team.Players {
		teamPlayersStars += player.Stars
	}

	switch action {
	case Contract:
		// Audit log embed

		var embed = discordgo.MessageEmbed{
			Title:     "✅ | **Recruitment Confirmed**",
			Timestamp: time.Now().String(),
			Color:     teamRole.Color,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: teamLogo,
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Creatine | Transactions",
				IconURL: "",
			},
			Description: fmt.Sprintf(""), //todo: whatever makes sense cause i forgo
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Player**",
					Value: fmt.Sprintf("- <@%s> \n - ⭐: %d", playerID, playerStars),
				},
				{
					Name:  "**Team**",
					Value: fmt.Sprintf("- <:%s:%s> **%s** \n - **Stars:** %d/%d \n - **Players:** %d/%d", teamSName, teamEID, teamLName, teamPlayersStars, teamSCap, teamPlayersNum, teamPCap),
				},
			},
		}

		// Send the embed

		response := discordgo.MessageSend{
			Embed: &embed,
		}

		_, err := s.ChannelMessageSendComplex(cfg.Transactions, &response)
		if err != nil {
			log.Println("error sending transaction")
			return
		}
	default:
		panic("unhandled default case")

	}

}
