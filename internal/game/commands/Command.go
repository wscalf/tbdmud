package commands

import (
	"github.com/wscalf/tbdmud/internal/game/commands/parameters"
	"github.com/wscalf/tbdmud/internal/game/world"
)

type Command interface {
	GetDescription() string
	GetParameters() []parameters.Parameter
	Execute(player *world.Player, args map[string]string)
}
