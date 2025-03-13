package game

import (
	"errors"
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game/commands"
	"github.com/wscalf/tbdmud/internal/game/contracts"
	"github.com/wscalf/tbdmud/internal/game/jobs"
	"github.com/wscalf/tbdmud/internal/game/world"
)

type Game struct {
	listener     contracts.ClientListener
	commands     *commands.Commands
	players      map[string]*world.Player
	jobQueue     *jobs.JobQueue
	startingRoom *world.Room
	login        *Login
}

func NewGame(commands *commands.Commands, listener contracts.ClientListener, startingRoom *world.Room, login *Login) *Game {
	return &Game{
		commands:     commands,
		listener:     listener,
		players:      make(map[string]*world.Player),
		jobQueue:     jobs.NewJobQueue(100),
		startingRoom: startingRoom,
		login:        login,
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
		p := g.login.Process(client)
		if p == nil {
			client.Disconnect()
			continue
		}

		g.players[p.ID] = p
		p.AttachClient(client)
		p.SetInputHandler(g.handleCommand)
		go func() {
			g.jobQueue.Enqueue(JoinWorldJob{
				player: p,
				game:   g,
			})
			p.Run()
			delete(g.players, p.ID)
		}()

	}

	err = g.listener.LastError()
	if err != nil {
		slog.Error("Error from client listener", "error", err)
	}
}

type JoinWorldJob struct {
	player *world.Player
	game   *Game
}

func (j JoinWorldJob) Run() {
	j.player.Join(j.game.startingRoom) //For persistence would need to retrieve the character's last room from storage
	j.player.SetInputHandler(j.game.handleCommand)
}

func (g *Game) handleCommand(player *world.Player, cmd string) {
	job, err := g.commands.Prepare(player, cmd)
	if err != nil {
		if errors.Is(err, commands.InputError) {
			player.Send(err.Error())
		} else {
			player.Send("An error has occurred.")
			slog.Error("Error processing command", "player", player.Name, "cmd", cmd, "err", err)
		}
		return
	}

	g.jobQueue.Enqueue(job)
}
