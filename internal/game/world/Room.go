package world

import (
	"github.com/wscalf/tbdmud/internal/game/contracts"
)

type Room struct {
	Object
	players map[string]*Player
}

func NewRoom(id string, name string, script contracts.ScriptObject) *Room {
	return &Room{
		Object: Object{
			ID:          id,
			Name:        name,
			Description: "",
			script:      script,
		},
		players: map[string]*Player{},
	}
}

func (r *Room) SendToAll(template string, params ...interface{}) {
	for _, p := range r.players {
		p.Send(template, params...)
	}
}

func (r *Room) SendToAllExcept(player *Player, template string, params ...interface{}) {
	for _, p := range r.players {
		if p == player {
			continue
		}

		p.Send(template, params...)
	}
}

func (r *Room) addPlayer(p *Player) {
	r.players[p.ID] = p

	p.Send("You joined: %s", r.Name)
}

func (r *Room) removePlayer(p *Player) {
	delete(r.players, p.ID)

	p.Send("You left: %s", r.Name)
}
