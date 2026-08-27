package game

import "log/slog"

type Players struct {
	idToPlayer map[string]*Player
	nameToId   map[string]string
	script     ScriptObject
}

func NewPlayers() *Players {
	return &Players{
		idToPlayer: map[string]*Player{},
		nameToId:   map[string]string{},
	}
}

func (p *Players) AttachScript(script ScriptObject) {
	p.script = script
}

func (p *Players) Add(player *Player) {
	p.nameToId[player.Name] = player.ID
	p.idToPlayer[player.ID] = player

	if p.script != nil {
		if _, err := p.script.Call("_on_player_joined", player); err != nil {
			slog.Error("script error processing player joined event", "player", player.Name, "error", err)
		}
	}
}

func (p *Players) Remove(player *Player) {
	if p.script != nil {
		if _, err := p.script.Call("_on_player_leaving", player); err != nil {
			slog.Error("script error processing player disconnecting event", "player", player.Name, "error", err)
		}
	}

	delete(p.idToPlayer, player.ID)
	delete(p.nameToId, player.Name)
}

func (p *Players) All() []*Player {
	count := len(p.idToPlayer)
	result := make([]*Player, 0, count)

	for _, player := range p.idToPlayer {
		result = append(result, player)
	}

	return result
}

func (p *Players) FindById(id string) *Player {
	return p.idToPlayer[id]
}

func (p *Players) FindByName(name string) *Player {
	id := p.nameToId[name]
	if id == "" {
		return nil
	}

	return p.idToPlayer[id]
}
