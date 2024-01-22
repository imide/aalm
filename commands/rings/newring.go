package rings

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/db"
	"log"
)

var NewRing = cmdutil.Commands{
	Name:        "newring",
	Description: "Creates a new ring.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "The FULL NAME of the ring.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "description",
			Description: "Provide a short and brief description of the ring.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role",
			Description: "The role to assign to the ring. If not provided, a role will automatically be created.",
			Required:    false,
		},
	},
	Handler:     newRingHandler,
	Permissions: discordgo.PermissionAdministrator,
}

func newRingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Button stuff (due to how i want the embed to work, only the buttons will  be here. Some embed stuff will be here too):

	var cancelEmbed = discordgo.MessageEmbed{
		Title:       "️❌| **Cancelled**",
		Description: "The ring creation has been cancelled.",
		Color:       0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL League Management",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
	}

	var successEmbed = discordgo.MessageEmbed{
		Title:       "✅ | **Success**",
		Description: fmt.Sprintf("The ring `**%s**` has been created. You may now assign this ring to players.", i.ApplicationCommandData().Options[0].StringValue()),
		Color:       0xffcc4d,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "AAFL League Management",
			IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
		},
	}

	var acceptButton = cmdutil.CreateButton("Accept", discordgo.SuccessButton, "✅", "accept")

	var denyButton = cmdutil.CreateButton("Deny", discordgo.DangerButton, "❌", "deny")

	var confirmRow = discordgo.ActionsRow{Components: []discordgo.MessageComponent{*acceptButton, *denyButton}}

	// Logic:

	if len(i.ApplicationCommandData().Options) > 2 {
		switch i.ApplicationCommandData().Options[2].StringValue() {
		default:
			// Warn not implemented
			cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "This feature is not implemented yet.", 0xffcc4d))
			//TODO: implement
		}
	} else {
		// Embed but also warn that a role will be created
		confirmEmbed := discordgo.MessageEmbed{
			Title:       "️⚠️ | **Warning**",
			Description: fmt.Sprintf("A ring will be created with the following data assigned:\n\n**Name:** %s\n**Description:** %s\n**Role:** None assigned, will automatically create one.\n\nAre you sure you want to continue?", i.ApplicationCommandData().Options[0].StringValue(), i.ApplicationCommandData().Options[1].StringValue()),
			Color:       0xffcc4d,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "AAFL League Management",
				IconURL: "https://cdn.discordapp.com/attachments/1182340550725214208/1182344030391111710/aafl_logo.png",
			},
		}

		// Send the embed with the buttons
		response := discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{&confirmEmbed},
				Components: []discordgo.MessageComponent{confirmRow},
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		}

		err := s.InteractionRespond(i.Interaction, &response)
		if err != nil {
			cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to your interaction.", 0xffcc4d))
			log.Println("Error responding to interaction,", err)
			return
		}

		// Handle the button presses

		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Type == discordgo.InteractionMessageComponent {
				switch i.ApplicationCommandData().ID {
				case "accept":
					// Create and edit the role so the Ring struct can be full
					roleParams := discordgo.RoleParams{
						Name: i.ApplicationCommandData().Options[0].StringValue(),
					}
					role, err := s.GuildRoleCreate(i.GuildID, &roleParams)
					if err != nil {
						cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while creating the role.", 0xffcc4d))
						log.Println("Error creating role,", err)
						return
					}

					// Create the ring struct
					ringData := db.Ring{
						Name:   i.ApplicationCommandData().Options[0].StringValue(),
						Desc:   i.ApplicationCommandData().Options[1].StringValue(),
						RoleID: role.ID,
					}

					// Create the ring
					err = db.CreateRing(ringData)
					if err != nil {
						cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while creating the ring.", 0xffcc4d))
						log.Println("Error creating ring,", err)
						return
					}

					// Send success embed

					response := discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseUpdateMessage,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{&successEmbed},
							Flags:  discordgo.MessageFlagsEphemeral,
						},
					}

					err = s.InteractionRespond(i.Interaction, &response)
					if err != nil {
						cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to your interaction.", 0xffcc4d))
						log.Println("Error responding to interaction,", err)
						return
					}

					//TODO: audit log

					return
				case "deny":
					// Send the embed

					response := discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseUpdateMessage,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{&cancelEmbed},
							Flags:  discordgo.MessageFlagsEphemeral,
						},
					}

					err = s.InteractionRespond(i.Interaction, &response)
					if err != nil {
						cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while responding to your interaction.", 0xffcc4d))
						log.Println("Error responding to interaction,", err)
						return
					}

					return
				}
			}
		})
	}
}
