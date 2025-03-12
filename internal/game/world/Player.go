package world

import (
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game/contracts"
)

type Player struct {
	Object
	client  contracts.Client
	onInput func(*Player, string)
}

func NewPlayer(id string, name string, script contracts.ScriptObject, client contracts.Client, onInput func(*Player, string)) *Player {
	return &Player{
		Object: Object{
			ID:     id,
			Name:   name,
			script: script,
		},
		client:  client,
		onInput: onInput,
	}
}

func (p *Player) Run() {
	for input := range p.client.Recv() {
		p.onInput(p, input)
	}

	err := p.client.LastError()
	if err != nil {
		slog.Error("Communication error from player", "name", p.Name, "error", err)
	}
	//Handle disconnect
}

func (p *Player) Send(text string) error {
	return p.client.Send(text)
}
