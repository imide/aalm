package admin

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/imide/aalm/commands/cmdutil"
	"github.com/imide/aalm/util/config"
	"github.com/imide/aalm/util/db"
	"log"
)

var AddTeam = cmdutil.Commands{
	Name:        "addteam",
	Description: "Add a team to the system",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "name",
			Description: "The FULL of the team",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "logo",
			Description: "The URL of the team's logo",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "emoji",
			Description: "The emoji ID of the team's emoji",
		},
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "owner",
			Description: "The owner of the team",
		},
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "role",
			Description: "The role of the team",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "invite",
			Description: "The discord invite of the team",
		},
	},
	Handler:     addTeamHandler,
	Permissions: discordgo.PermissionAdministrator,
}

func addTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extracting options
	options := i.ApplicationCommandData().Options
	name := options[0].StringValue()
	logo := options[1].StringValue()
	emoji := options[2].StringValue()
	owner := options[3].UserValue(s)
	role := options[4].RoleValue(s, i.GuildID)
	invite := options[5].StringValue()

	// Creating the team
	newTeam := &db.Team{
		Name:          name,
		Logo:          logo,
		EmojiID:       emoji,
		Owner:         owner.ID,
		RoleID:        role.ID,
		DiscordInvite: invite,
	}

	// Adding the team to the database
	err := db.SaveTeamData(newTeam)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while saving the team data.", 0xffcc4d))
		log.Println("Error saving team data,", err)
		return
	}

	// Adding the team role to the owner
	err = s.GuildMemberRoleAdd(i.GuildID, owner.ID, role.ID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while adding the team role to the owner.", 0xffcc4d))
		log.Println("Error adding team role to owner,", err)
		return
	}

	// Adding the FO role to the owner
	err = s.GuildMemberRoleAdd(i.GuildID, owner.ID, config.Cfg.FORoleID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while adding the FO role to the owner.", 0xffcc4d))
		log.Println("Error adding FO role to owner,", err)
		return
	}

	// DMing the owner
	_, err = s.UserChannelCreate(owner.ID)
	if err != nil {
		cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("⚠️ | **Warning**", "An error occurred while creating the DM channel.", 0xffcc4d))
		log.Println("Error creating DM channel,", err)
		return
	}

	// Sending the DM
	_, err = s.ChannelMessageSendEmbed(owner.ID, cmdutil.CreateEmbed("✅ | **Success**", fmt.Sprintf("Your new team, %s, has been added to the database and fully registered. You now have full access to manage your team.", newTeam.Name), 0x00ff00))

	// Sending the confirmation message
	cmdutil.SendInteractionResponse(s, i, cmdutil.CreateEmbed("✅ | **Success**", "Successfully added the team to the database. The owner now has their team role and FO role.", 0x00ff00))

	//TODO: audit log
}
