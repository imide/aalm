package player

import "github.com/imide/aalm/commands"

var invalidRecruiterPerm = commands.CreateEmbed(
	"⚠️ | **Warning**",
	"You do not have permission to recruit players.\n\n"+"You must be either the team owner or a coach to recruit players.",
	0xffcc4d,
)

var recruitSuspendedErr = commands.CreateEmbed(
	"⚠️ | **Warning**",
	"Either you or the recruit is suspended.\n\n"+"Please wait until their suspension expires before recruiting them.",
	0xffcc4d,
)

var recruitAlreadyContractedErr = commands.CreateEmbed(
	"⚠️ | **Warning**",
	"The user you are trying to recruit is already contracted.\n\n"+"Please wait until their contract expires or they are dropped before recruiting them.",
	0xffcc4d,
)

var recruitOfferReceived = commands.CreateEmbed(
	"⚠️ | **Offer Received**",
	"The team (add here) has offered you a spot on their team within the AAFL. \nThis offer is optional; you are not inclined to accept the position.\n\n**This offer will expire in 24 hours.**",
	0x00ff00,
)

var offerExpired = commands.CreateEmbed(
	"⚠️ | **Offer Expired**",
	"The offer to join the team has expired.",
	0xffcc4d,
)
