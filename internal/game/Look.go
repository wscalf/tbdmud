package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

var lookparams []parameters.Parameter = []parameters.Parameter{parameters.NewName("object", false)}

type Look struct{}

func (l Look) GetDescription() string {
	return "Describes the current room or the object looked at."
}

func (l Look) GetParameters() []parameters.Parameter {
	return lookparams
}

func (l Look) Execute(player *Player, args map[string]string) {
	name, found := args["object"]
	room := player.GetRoom()

	if !found {
		player.Send(room.Describe())
		return
	}

	other := room.FindPlayer(name)
	if other != nil {
		player.Send(other.Describe())
		other.Sendf("%s looked at you", player.Name)
		return
	}

	player.Sendf("I don't see that here.")
}
