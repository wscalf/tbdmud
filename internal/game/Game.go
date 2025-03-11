package game

import (
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game/contracts"
)

type Game struct {
	listener contracts.ClientListener
	commands *Commands
	players  []*Player
}

func NewGame(commands *Commands, listener contracts.ClientListener) *Game {
	return &Game{
		commands: commands,
		listener: listener,
		players:  []*Player{},
	}
}

func (g *Game) Run() {
	g.handlePlayersJoining()
}

func (g *Game) Stop() {

}

func (g *Game) handlePlayersJoining() {
	clients, err := g.listener.Listen()
	if err != nil {
		slog.Error("Error listening for new clients", "error", err)
		return //Report this back up somehow?
	}

	for client := range clients {
		p := NewPlayer("", "", nil, client, g.commands)

		g.players = append(g.players, p)
		go p.Run()
	}

	err = g.listener.LastError()
	if err != nil {
		slog.Error("Error from client listener", "error", err)
	}
}
