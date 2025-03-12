package builtins

import (
	"github.com/wscalf/tbdmud/internal/game/commands/parameters"
	"github.com/wscalf/tbdmud/internal/game/world"
)

var params = []parameters.Parameter{parameters.NewFreeText("thought")}

type Think struct {
}

func (t Think) GetDescription() string {
	return "Sends the thought text to the player."
}

func (t Think) GetParameters() []parameters.Parameter {
	return params
}

func (t Think) Execute(player *world.Player, args map[string]string) {
	thought := args["thought"]

	player.Send(thought)
}
