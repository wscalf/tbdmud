package game

import "github.com/wscalf/tbdmud/internal/text"

type Room struct {
	Object
	players map[string]*Player
	links   map[string]*Link
	layout  *text.Layout
}

func NewRoom(id string, name string, description string, script ScriptObject) *Room {
	return &Room{
		Object: Object{
			ID:          id,
			Name:        name,
			Description: description,
			script:      script,
		},
		players: map[string]*Player{},
		links:   map[string]*Link{},
	}
}

func (r *Room) Link(command string, name string, description string, to *Room) {
	link := &Link{
		Object: Object{
			ID:          r.ID + "_" + command,
			Name:        name,
			Description: description,
			script:      nil,
		},
		command: command,
		to:      to,
	}

	r.links[command] = link
}

func (r *Room) FindLocalCommand(command string) Command {
	if link, ok := r.links[command]; ok {
		return link
	}

	return nil
}

func (r *Room) Describe() text.FormatJob {
	return r.layout.Prepare(r)
}

func (r *Room) GetProperties() map[string]interface{} {
	props := map[string]interface{}{
		"name": r.Name,
		"desc": r.Description,
	}

	players := make([]map[string]interface{}, 0, len(r.players))
	for _, p := range r.players {
		players = append(players, p.GetProperties())
	}
	props["players"] = players

	links := make([]map[string]interface{}, 0, len(r.links))
	for _, l := range r.links {
		links = append(links, l.GetProperties())
	}
	props["links"] = links

	return props
}

func (r *Room) SendToAll(template string, params ...interface{}) {
	for _, p := range r.players {
		p.Sendf(template, params...)
	}
}

func (r *Room) SendToAllExcept(player *Player, template string, params ...interface{}) {
	for _, p := range r.players {
		if p == player {
			continue
		}

		p.Sendf(template, params...)
	}
}

func (r *Room) addPlayer(p *Player) {
	r.players[p.ID] = p

	p.Send(r.Describe())
}

func (r *Room) removePlayer(p *Player) {
	delete(r.players, p.ID)

	p.Sendf("You left: %s", r.Name)
}
