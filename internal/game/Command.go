package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

type Command interface {
	GetDescription() string
	GetParameters() []parameters.Parameter
	Execute(player *Player, args map[string]string, state map[string]any, requeueHandler func()) bool
}
