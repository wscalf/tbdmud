package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

type Link struct {
	Object
	command string
	to      *Room
}

func (l *Link) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"name": l.Name,
		"desc": l.Description,
		"cmd":  l.command,
	}
}

// Run implements jobs.Job.
func (l *Link) GetDescription() string {
	return l.Description //Probably not used
}
func (l *Link) GetParameters() []parameters.Parameter {
	return []parameters.Parameter{} //Links are always parameterless
}

func (l *Link) Move(player *Player, to *Room) {
	if player.room != nil {
		player.Leave()
	}

	player.Join(l.to)
}

func (l *Link) Execute(player *Player, args map[string]string) {
	l.script.Call("Move", player, l.to)
}
