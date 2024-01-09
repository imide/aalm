package auditlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm"
	"github.com/imide/aalm/util/db"
	"time"
)

// This will contain the audit log functions relating to all transactions of players (trades, cuts, suspensions, etc.)

type TransactionType int

const (
	Trade TransactionType = iota
	Cut
	Suspend
	Unsuspend
	Contract
	Drop
)

//TODO: add input params cause ik i need more im just doing contracting for now

func LogTransaction(s *discordgo.Session, action TransactionType, player db.Player, team db.Team) {
	var cfg = main.Cfg

	switch action {
	case Contract:
		// Variable shortcuts cause im fucking lazy and i dont care
		var playerID, playerStars = player.ID, player.Stars // player shortcuts (add when needed)
		var teamRID, teamName, teamPCap, teamSCap, teamLogo = team.RoleID, team.Name, team.PlayerMax, team.MaxStars, team.Logo
		var teamRole, _ = s.State.Role(cfg.GuildID, teamRID)

		// Audit log embed

		var embed = discordgo.MessageEmbed{
			Title:     "✅ | **Transaction Confirmed**",
			Timestamp: time.Now().String(),
			Color:     teamRole.Color,
		}

	}

}
