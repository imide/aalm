package auditlog

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/util/db"
	"time"
)

type SuspensionAction int

const (
	Suspend SuspensionAction = iota
	Unsuspend
)

var ErrNotImplemented = errors.New("not implemented")

// LogSuspension logs a suspension to the audit log (placeholder)
func LogSuspension(s *discordgo.Session, player db.Player, action SuspensionAction) error {
	return ErrNotImplemented
}
