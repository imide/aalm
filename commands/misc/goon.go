package misc

import "github.com/imide/aalm/commands"

var goon = commands.Commands{
	Name:        "goon",
	Description: "Goon.",
	Options:     nil,
	Handler:     goonHandler,
}
