package team

import (
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/db"
)

var AddCoach = cmdutil.Commands{
	Name:        "addcoach",
	Description: "Adds a coach to the team.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "coach",
			Description: "The coach to add.",
			Required:    true,
		},
	},
	Handler: addCoachHandler,
}

func addCoachHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get important data
	ownerData, err := db.GetPlayerData(i.Member.User.ID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving your data.", 0xffcc4d))
		return
	}

	teamData, err := db.GetTeamData(ownerData.TeamPlaying)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving your team data.", 0xffcc4d))
		return
	}

	coachData, err := db.GetPlayerData(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while retrieving the coach data.", 0xffcc4d))
		return
	}

	// Permission check
	if teamData.Owner != ownerData.ID {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "You do not have permission to use this command.", 0xffcc4d))
		return
	}

	// Other misc checks

	switch {
	case coachData.TeamPlaying != teamData.ID:
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "This user is not on your team.", 0xffcc4d))
		return
	case coachData.Suspension != nil && coachData.Suspension.Active == true:
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "This user is suspended.", 0xffcc4d))
		return
	case teamData.Coaches != nil && len(teamData.Coaches) >= 3:
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "You cannot have more than 3 coaches.", 0xffcc4d))
	case slices.Contains(teamData.Coaches, coachData.ID):
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "This user is already a coach on your team.", 0xffcc4d))
		return
	// TODO: whatever is missing
	default:
		// Add the coach to the team
		teamData.Coaches = append(teamData.Coaches, coachData.ID)

		// Add the role to the coach
		// TODO: add coach role to config

		// Save
		err := db.SaveTeamData(&teamData)
		if err != nil {
			cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while saving the team data.", 0xffcc4d))
			return
		}
	}
}
