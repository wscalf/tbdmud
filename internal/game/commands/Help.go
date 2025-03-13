package commands

import (
	"github.com/wscalf/tbdmud/internal/game/commands/parameters"
	"github.com/wscalf/tbdmud/internal/game/world"
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

func (h Help) Execute(player *world.Player, params map[string]string) {
	if name, ok := params["cmd"]; ok {
		//A specific command was passed
		if cmd, ok := h.commands.commands[name]; ok {
			params := cmd.GetParameters()
			player.Send("%s: %s", name, cmd.GetDescription())
			usage := "Usage: " + name
			for _, param := range params {
				usage = usage + " [" + param.Name() + "]"
			}
			player.Send(usage)
		} else {
			player.Send("Unrecognized command: %s", name)
		}
	} else {
		//No command was passed, print the list
		//This should probably be externalized to a template
		player.Send("The following commands are available:")
		for name, cmd := range h.commands.commands {
			player.Send("%s: %s", name, cmd.GetDescription())
		}
	}
}
