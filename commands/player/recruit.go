package player

import (
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands"
	"github.com/imide/aalm/util/db"
	"log"
)

var recruit = commands.Commands{
	Name:        "recruit",
	Description: "Recruit a user to the group.",
	Options:     nil,
	Handler:     recruitHandler,
}

func recruitHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Fetch the user's role and the player's data
	recruiterData, err := db.GetPlayerData(i.Member.User.ID)
	recruitId := i.ApplicationCommandData().Options[0].UserValue(s).ID
	recruitUserData, err := db.GetPlayerData(recruitId)
	if err != nil {
		// Handle error
		return
	}

	validPermissions := []string{"owner", "coach"}

	switch {
	case !contains(validPermissions, recruiterData.Position):
		commands.SendInteractionResponse(s, i, invalidRecruiterPerm)
		log.Println("interaction recruit denied, invalid recruiter permissions")
		return
	case recruitUserData.IsSuspended == true || recruiterData.IsSuspended == true:
		commands.SendInteractionResponse(s, i, recruitSuspendedErr)
		log.Println("interaction recruit denied, recruit/recruiter is suspended")
		return
	case recruitUserData.Contracted == true || recruitUserData.TeamPlaying != "":
		commands.SendInteractionResponse(s, i, recruitAlreadyContractedErr)
		log.Println("interaction recruit denied, recruit is already contracted")
		return
	default:
		confirm := commands.CreateButton("Confirm", discordgo.SuccessButton, "✅", "confirm")
		deny := commands.CreateButton("Deny", discordgo.DangerButton, "❌", "deny")
		components := discordgo.ActionsRow{Components: []discordgo.MessageComponent{*confirm, *deny}}
		commands.SendInteractionWithComponent(s, i, recruitOfferReceived, []discordgo.MessageComponent{components})

		// If no response is received within 24 hours, automatically deny the recruitment
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
