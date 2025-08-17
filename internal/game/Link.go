package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

type Link struct {
	Object
	Command string
	to      *Room
	from    *Room
}

func (l *Link) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"name": l.Name,
		"desc": l.Description,
		"cmd":  l.Command,
	}
}

// Run implements jobs.Job.
func (l *Link) GetDescription() string {
	return l.Description //Probably not used
}
func (l *Link) GetParameters() []parameters.Parameter {
	return []parameters.Parameter{} //Links are always parameterless
}

func (l *Link) Peek() *Room {
	return l.to
}

func (l *Link) Move(player *Player, to *Room) {
	if player.room != nil {
		player.Leave()
	}

	player.Join(l.to)
}

func (l *Link) Execute(player *Player, args map[string]string, state map[string]any, requeueHandler func()) bool {
	l.script.Call("Move", player, l.to)
	return true
}
