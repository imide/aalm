package player

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/auditlog"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/db"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type RecruitmentOffer struct {
	UserID   string
	TimeSent time.Time
}

var Recruit = cmdutil.Commands{
	Name:        "recruit",
	Description: "Recruit a player to your team.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "player",
			Description: "The player to recruit.",
			Required:    true,
		},
	},
	Handler: recruitHandler,
}

// Functions:

func recruitHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Variables and checking:

	recruiterData, recruiterTeam, err := GetRecruiterData(s, i)

	if err != nil {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while trying to get your team data. Please try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return
	}

	guildID := i.GuildID

	recruitData, isRecruitable, err := getRecruitData(s, i)
	switch isRecruitable {
	case false:
		return
	}

	switch err {
	case nil:
		break
	default:
		return
	}

	teamRole, err := s.State.Role(guildID, recruiterTeam.RoleID)
	if err != nil {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while trying to get your team data. Please try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return
	}

	// Messages and components (well important ones):

	var recruitDm = discordgo.MessageEmbed{
		Title: "📨 | **Recruitment Offer**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`The coach of the team **%s** has offered you a contract to play for their team within the AAFL. Would you like to accept? \n \n **Note:** You will be unable to play for another team until your contract expires. **This offer will expire in 24 hours.**`, recruiterTeam.Name),
		Color:       teamRole.Color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: recruiterTeam.Logo,
		},
	}

	var recruitAccept = discordgo.MessageEmbed{
		Title: "✅ | **Recruitment Accepted**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`Welcome to the **%s**. Your recruitment has been finalized. \n You may join your team's discord [here.](https://discord.gg/%s)\n **You will be unable to play for another team until your contract expires.**`, recruiterTeam.Name, recruiterTeam.DiscordInvite),
		Color:       teamRole.Color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: recruiterTeam.Logo,
		},
	}

	var recruitDeny = discordgo.MessageEmbed{
		Title: "❌ | **Recruitment Denied**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`You have denied the recruitment offer from the **%s**.`, recruiterTeam.Name),
		Color:       teamRole.Color,
	}

	var recruitTimeout = discordgo.MessageEmbed{
		Title: "⏰ | **Recruitment Expired**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`The recruitment offer from the **%s** has expired.`, recruiterTeam.Name),
		Color:       teamRole.Color,
	}

	var recruitmentSent = discordgo.MessageEmbed{
		Title: "📨 | **Recruitment Sent**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`Your recruitment offer has been sent to <@%s>.`, recruitData.ID),
	}

	var recruitmentDenied = discordgo.MessageEmbed{
		Title: "❌ | **Recruitment Denied**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`Your recruitment offer to <@%s> was denied or ignored by the user.`, recruitData.ID),
	}

	var recruitmentAccepted = discordgo.MessageEmbed{
		Title: "✅ | **Recruitment Accepted**",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL Recruitment",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
		Description: fmt.Sprintf(`Your recruitment offer to <@%s> was accepted. They have recieved their roles and their data has been updated as well`, recruitData.ID),
	}

	var acceptButton = cmdutil.CreateButton("Accept", discordgo.SuccessButton, "✅", "accept")

	var denyButton = cmdutil.CreateButton("Deny", discordgo.DangerButton, "❌", "deny")

	var confirmRow = discordgo.ActionsRow{Components: []discordgo.MessageComponent{*acceptButton, *denyButton}}

	// Actual code:

	var recruitmentOffer = RecruitmentOffer{
		UserID:   recruitData.ID,
		TimeSent: time.Now(),
	}

	cmdutil.SendInteractionResponse(s, i, &recruitmentSent)

	dmChannel, err := s.UserChannelCreate(recruitData.ID)
	if err != nil {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while trying to send the recruitment offer. Please try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return
	}

	message := &discordgo.MessageSend{
		Embed:      &recruitDm,
		Components: []discordgo.MessageComponent{confirmRow},
	}

	_, err = s.ChannelMessageSendComplex(dmChannel.ID, message)
	if err != nil {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while trying to send the recruitment offer. Please try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return
	}

	// Add a handler for InteractionCreate events
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Check if the interaction is a MessageComponent interaction
		if i.Type == discordgo.InteractionMessageComponent {
			// Check the CustomID of the button
			switch i.MessageComponentData().CustomID {
			case "accept":
				// Handle the "Accept" button

				if time.Now().Sub(recruitmentOffer.TimeSent).Hours() > 24 {
					// The offer has expired
					messageEdit := &discordgo.MessageEdit{
						ID:      i.Message.ID,
						Channel: i.ChannelID,
						Embed:   &recruitTimeout,
					}
					_, err := s.ChannelMessageEditComplex(messageEdit)
					if err != nil {
						log.Println("Error editing message,", err)
					}
					return
				}
				// Edit the message to show the recruitAccept embed
				messageEdit := &discordgo.MessageEdit{
					ID:      i.Message.ID,
					Channel: i.ChannelID,
					Embed:   &recruitAccept,
				}
				_, err := s.ChannelMessageEditComplex(messageEdit)
				if err != nil {
					log.Println("Error editing message,", err)
				}

				// Update the player's data

				recruitData.TeamPlaying = recruiterTeam.ID
				recruitData.Contracted = true
				err = db.SavePlayerData(&recruitData)

				// Create a DM channel with the recruiter
				dmChannel, err := s.UserChannelCreate(recruiterData.ID)
				if err != nil {
					log.Println("Error creating DM channel,", err)
					return
				}

				// Send a message to the recruiter in their DMs
				_, err = s.ChannelMessageSendEmbed(dmChannel.ID, &recruitmentAccepted)
				if err != nil {
					log.Println("Error sending message to DM channel,", err)
				}

				// Log the transaction
				auditlog.LogTransaction(s, auditlog.Contract, recruitData, recruiterTeam)

			case "decline":
				// Handle the "Decline" button
				messageEdit := &discordgo.MessageEdit{
					ID:      i.Message.ID,
					Channel: i.ChannelID,
					Embed:   &recruitDeny,
				}
				_, err := s.ChannelMessageEditComplex(messageEdit)
				if err != nil {
					log.Println("Error editing message,", err)
				}

				// Create a DM channel with the recruiter
				dmChannel, err := s.UserChannelCreate(recruiterData.ID)
				if err != nil {
					log.Println("Error creating DM channel,", err)
					return
				}

				// Send a message to the recruiter in their DMs
				_, err = s.ChannelMessageSendEmbed(dmChannel.ID, &recruitmentDenied)
				if err != nil {
					log.Println("Error sending message to DM channel,", err)
				}
			}
		}
	})

}

func GetRecruiterData(s *discordgo.Session, i *discordgo.InteractionCreate) (db.Player, db.Team, error) {
	recruiterData, err := db.GetPlayerData(i.Member.User.ID)
	if err != nil {
		return db.Player{}, db.Team{}, err
	}

	recruiterTeam, _ := db.GetTeamData(recruiterData.TeamPlaying)

	if !db.HasManagePermission(i.Member.User.ID, recruiterTeam) {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "You do not have permission to manage this team. Either run setteam or try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		noperm := errors.New("no permission")
		return db.Player{}, db.Team{}, noperm
	}

	return recruiterData, recruiterTeam, nil
}

func getRecruitData(s *discordgo.Session, i *discordgo.InteractionCreate) (db.Player, bool, error) {
	recruitData, err := db.GetPlayerData(i.ApplicationCommandData().Options[0].UserValue(s).ID)
	if err != nil {
		switch err.Error() {
		case mongo.ErrNoDocuments.Error():
			embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "The user specified is not registered in the database. Creating a user for you...", 0xffcc4d)
			cmdutil.SendInteractionResponse(s, i, embed)
			player, err := cmdutil.UserDoesntExist(i.ApplicationCommandData().Options[0].UserValue(s).ID)
			if err != nil {
				cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An unknown error occurred.", 0xffcc4d))
				log.Printf("an error occurred during getting data: %s", err)
				return db.Player{}, false, err
			}

			return player, true, nil

		default:
			return db.Player{}, false, err
		}
	}

	if recruitData.TeamPlaying != "" || recruitData.Contracted == true {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", fmt.Sprintf(`The player you are trying to recruit is already playing for the **%s**. \n\n You may trade for the player or try again later.`, recruitData.TeamPlaying), 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return db.Player{}, false, nil
	}

	if recruitData.Suspension != nil && recruitData.Suspension.Active == true {
		embed := cmdutil.CreateEmbed("⚠️ | **Warning**", "This player is currently suspended. Please try again later.", 0xffcc4d)
		cmdutil.SendInteractionResponse(s, i, embed)
		return db.Player{}, false, nil
	}

	return recruitData, true, nil
}
