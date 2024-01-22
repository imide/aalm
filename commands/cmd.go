package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/admin"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/commands/misc"
	"github.com/imide/aalm/commands/player"
	"github.com/imide/aalm/commands/rings"
	"github.com/imide/aalm/commands/team"
	"log"
)

var CmdMap = map[string]cmdutil.Commands{}

func init() {
	// team
	CmdMap["roster"] = team.Roster
	CmdMap["release"] = team.Release
	// rings
	CmdMap["rings"] = rings.Rings
	CmdMap["newring"] = rings.NewRing
	// player
	CmdMap["ringcheck"] = player.Ringcheck
	CmdMap["recruit"] = player.Recruit
	CmdMap["drop"] = player.Drop
	// misc
	CmdMap["goon"] = misc.Goon
	// admin
	CmdMap["suspend"] = admin.Suspend
	CmdMap["setteam"] = admin.SetTeam
}

func Register(s *discordgo.Session, guildID string) {
	// Fetch existing commands
	existingCommands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		log.Fatalf("Failed to fetch existing commands: %v", err)
	}

	// Create a map of local commands for easy lookup
	localCommands := make(map[string]cmdutil.Commands)
	for _, cmd := range CmdMap {
		localCommands[cmd.Name] = cmd
	}

	// Delete commands that exist on Discord but not locally
	for _, cmd := range existingCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
		if err != nil {
			log.Println("Err deleting cmd,", err)
		}
	}

	// Register local commands
	for _, command := range CmdMap {
		createCmd(s, guildID, command)
	}
}

func createCmd(s *discordgo.Session, guildID string, cmd cmdutil.Commands) {
	perm := int64(cmd.Permissions)

	_, err := s.ApplicationCommandCreate(s.State.Application.ID, guildID, &discordgo.ApplicationCommand{
		Name:                     cmd.Name,
		Description:              cmd.Description,
		Options:                  cmd.Options,
		DefaultMemberPermissions: &perm,
	})
	if err != nil {
		log.Println("Err creating cmd,", err)
	}
}
