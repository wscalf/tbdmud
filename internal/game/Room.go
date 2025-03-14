package game

type Room struct {
	Object
	players map[string]*Player
	links   map[string]*Link
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
