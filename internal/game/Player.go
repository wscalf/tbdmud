package game

import (
	"fmt"
	"log/slog"
)

type Player struct {
	Object
	client  Client
	room    *Room
	outbox  chan output
	onInput func(*Player, string)
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		Object: Object{
			ID:   id,
			Name: name,
		},
		outbox: make(chan output, 10),
	}
}

func (p *Player) AttachScript(script ScriptObject) {
	p.script = script
}

func (p *Player) AttachClient(client Client) {
	p.client = client
}

func (p *Player) SetInputHandler(onInput func(*Player, string)) {
	p.onInput = onInput
}

func (p *Player) GetRoom() *Room {
	return p.room
}

func (p *Player) Join(r *Room) {
	r.addPlayer(p)
	p.room = r
}

func (p *Player) Leave() {
	p.room.removePlayer(p)
	p.room = nil
}

func (p *Player) Run() {
	inbox := p.client.Recv()
	active := true

	for active {
		select {
		case input, ok := <-inbox:
			if ok {
				p.onInput(p, input)
			} else {
				err := p.client.LastError()
				if err != nil {
					slog.Error("Communication error from player", "name", p.Name, "error", err)
				}
				active = false
			}
		case output := <-p.outbox:
			msg := fmt.Sprintf(output.template, output.params...)
			p.client.Send(msg)
		}
	}
}

func (p *Player) Send(template string, params ...interface{}) {
	p.outbox <- output{template: template, params: params}
}

type output struct {
	template string
	params   []interface{}
}
