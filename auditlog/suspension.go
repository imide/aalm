package auditlog

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/util/config"
	"github.com/imide/aalm/util/db"
	"time"
)

type SuspensionAction int

const (
	Suspend SuspensionAction = iota
	Unsuspend
)

func LogSuspension(s *discordgo.Session, player db.Player, action SuspensionAction) error {
	var cfg = config.Cfg
	var suspension = player.Suspension // doing this because fuck it

	// whatever else is needed like vars and shit goes above you fucking retard
	switch action {

	case Suspend:
		var embed = &discordgo.MessageEmbed{
			Title:     "✅ | **Suspension Confirmed**",
			Timestamp: time.Now().String(),
			Color:     0x800000,
		}
	}
}
