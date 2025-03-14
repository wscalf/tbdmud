package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

var sayparams []parameters.Parameter = []parameters.Parameter{parameters.NewFreeText("text")}

type Say struct{}

func (s Say) GetDescription() string {
	return "Sends the thought text to the player."
}

func (s Say) GetParameters() []parameters.Parameter {
	return sayparams
}

func (s Say) Execute(player *Player, args map[string]string) {
	text := args["text"]
	player.Send(`You say, "%s"`, text)

	room := player.GetRoom()
	room.SendToAllExcept(player, `%s says, "%s"`, player.Name, text)
}
