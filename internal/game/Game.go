package game

import (
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game/commands"
	"github.com/wscalf/tbdmud/internal/game/contracts"
	"github.com/wscalf/tbdmud/internal/game/jobs"
	"github.com/wscalf/tbdmud/internal/game/world"
)

type Game struct {
	listener contracts.ClientListener
	commands *commands.Commands
	players  []*world.Player
	jobQueue *jobs.JobQueue
}

func NewGame(commands *commands.Commands, listener contracts.ClientListener) *Game {
	return &Game{
		commands: commands,
		listener: listener,
		players:  []*world.Player{},
		jobQueue: jobs.NewJobQueue(100),
	}
}

func (g *Game) Run() {
	go g.jobQueue.Run()
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
		p := world.NewPlayer("", "", nil, client, g.handleCommand)

		g.players = append(g.players, p)
		go p.Run()
	}

	err = g.listener.LastError()
	if err != nil {
		slog.Error("Error from client listener", "error", err)
	}
}

func (g *Game) handleCommand(player *world.Player, cmd string) {
	job := g.commands.Prepare(player, cmd)

	g.jobQueue.Enqueue(job)
}
