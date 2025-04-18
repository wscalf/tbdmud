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
	items   map[string]*Object
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
		items:  map[string]*Object{},
		outbox: make(chan text.FormatJob, 10),
	}
}

type PlayerSaveData struct {
	RoomID string `json:"room"`
	ObjectSaveData
}

func PlayerFromSaveData(data map[string]any) (*Player, error) {
	//Need to look up an maybe apply the past room
	p := &Player{
		Object: ObjectFromSaveData(data),
		items:  map[string]*Object{},
		outbox: make(chan text.FormatJob, 10),
	}

	if script, err := _scriptSystem.Wrap(p, data["type"].(string)); err != nil { //This needs to be on the main thread too, actually
		return nil, err
	} else {
		p.script = script
		//Apply saved properties ... on the main thread
	}

	if roomId, ok := data["room"]; ok {
		_ = _world.FindRoom(roomId.(string))
		//How to enqueue joining the room?
	}

	return p, nil
}

func (p *Player) GetSaveData() (*PlayerSaveData, error) {
	obj, err := p.Object.GetSaveData()
	if err != nil {
		return nil, err
	}

	data := &PlayerSaveData{
		ObjectSaveData: obj,
	}

	if p.room != nil {
		data.RoomID = p.room.ID
	}

	return data, nil
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
	props := p.Object.GetProperties()

	objects := make([]map[string]interface{}, 0, len(p.items))
	for _, o := range p.items {
		objects = append(objects, o.GetProperties())
	}
	props["items"] = objects

	return props
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

func (p *Player) FindItem(item string) *Object {
	return p.items[item]
}

func (p *Player) Give(item *Object) {
	//TODO: emit some kind of inventory-disturbed script event here
	p.items[item.Name] = item
}

func (p *Player) Take(item *Object) {
	//TODO: emit script event here
	delete(p.items, item.Name)
}

func (p *Player) GetItems() []*Object {
	items := make([]*Object, 0, len(p.items))
	for _, item := range p.items {
		items = append(items, item)
	}

	return items
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
