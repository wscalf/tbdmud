package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

var helpparams []parameters.Parameter = []parameters.Parameter{parameters.NewName("cmd", false)}

type Help struct {
	commands *Commands
}

func (h Help) GetDescription() string {
	return "Prints this help text."
}

func (h Help) GetParameters() []parameters.Parameter {
	return helpparams
}

func (h Help) Execute(player *Player, params map[string]string) {
	if name, ok := params["cmd"]; ok {
		//A specific command was passed
		if cmd, ok := h.commands.commands[name]; ok {
			params := cmd.GetParameters()
			player.Sendf("%s: %s", name, cmd.GetDescription())
			usage := "Usage: " + name
			for _, param := range params {
				usage = usage + " [" + param.Name() + "]"
			}
			player.Sendf(usage)
		} else {
			player.Sendf("Unrecognized command: %s", name)
		}
	} else {
		//No command was passed, print the list
		//This should probably be externalized to a template
		player.Sendf("The following commands are available:")
		for name, cmd := range h.commands.commands {
			player.Sendf("%s: %s", name, cmd.GetDescription())
		}
	}
}
