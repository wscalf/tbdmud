package game

import (
	"fmt"
	"log/slog"

	"github.com/wscalf/tbdmud/internal/text"
)

type Player struct {
	Object
	client  Client
	room    *Room
	outbox  chan text.FormatJob
	layout  *text.Layout
	onInput func(*Player, string)
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		Object: Object{
			ID:   id,
			Name: name,
		},
		outbox: make(chan text.FormatJob, 10),
	}
}

func (p *Player) AttachScript(script ScriptObject) {
	p.script = script
}

func (p *Player) AttachClient(client Client) {
	p.client = client
}

func (p *Player) SetLayout(layout *text.Layout) {
	p.layout = layout
}

func (p *Player) SetInputHandler(onInput func(*Player, string)) {
	p.onInput = onInput
}

func (p *Player) Describe() text.FormatJob {
	return p.layout.Prepare(p)
}

func (p *Player) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"name": p.Name,
		"desc": p.Description,
	}
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
			msg, err := output.Run()
			if err != nil {
				slog.Error("error formatting output", "err", err, "job", output)
			} else {
				p.client.Send(msg)
			}
		}
	}
}

func (p *Player) Sendf(template string, params ...interface{}) {
	p.outbox <- sprintfJob{template: template, params: params}
}

func (p *Player) Send(job text.FormatJob) {
	p.outbox <- job
}

type sprintfJob struct {
	template string
	params   []interface{}
}

func (f sprintfJob) Run() (string, error) {
	return fmt.Sprintf(f.template, f.params...), nil
}
