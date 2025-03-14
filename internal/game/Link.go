package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

type Link struct {
	Object
	command string
	to      *Room
}

// Run implements jobs.Job.
func (l *Link) GetDescription() string {
	return l.Description //Probably not used
}
func (l *Link) GetParameters() []parameters.Parameter {
	return []parameters.Parameter{} //Links are always parameterless
}
func (l *Link) Execute(player *Player, args map[string]string) {
	if player.room != nil {
		player.Leave()
	}

	player.Join(l.to)
}
