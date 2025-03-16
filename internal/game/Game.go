package game

import (
	"errors"
	"log/slog"

	"github.com/wscalf/tbdmud/internal/text"
)

type Game struct {
	listener          ClientListener
	commands          *Commands
	players           map[string]*Player
	jobQueue          *JobQueue
	world             *World
	login             *Login
	layouts           map[string]*text.Layout
	scriptSystem      ScriptSystem
	defaultPlayerType string
}

func NewGame(commands *Commands, listener ClientListener, world *World, login *Login, layouts map[string]*text.Layout, scriptSystem ScriptSystem, defaultPlayerType string) *Game {
	return &Game{
		commands:          commands,
		listener:          listener,
		players:           make(map[string]*Player),
		jobQueue:          NewJobQueue(100),
		world:             world,
		login:             login,
		layouts:           layouts,
		scriptSystem:      scriptSystem,
		defaultPlayerType: defaultPlayerType,
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
		go func() {
			p := g.login.Process(client)
			if p == nil {
				client.Disconnect()
				return
			}

			g.players[p.ID] = p
			p.AttachClient(client)
			p.SetInputHandler(g.handleCommand)
			script, _ := g.scriptSystem.Wrap(p, g.defaultPlayerType)
			p.AttachScript(script)

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
	player *Player
	game   *Game
}

func (j JoinWorldJob) Run() {
	j.player.Join(j.game.world.chargen) //For persistence would need to retrieve the character's last room from storage
	j.player.SetInputHandler(j.game.handleCommand)
}

func (g *Game) handleCommand(player *Player, cmd string) {
	job, err := g.commands.Prepare(player, cmd)
	if err != nil {
		if errors.Is(err, ErrInputError) {
			player.Sendf(err.Error())
		} else {
			player.Sendf("An error has occurred.")
			slog.Error("Error processing command", "player", player.Name, "cmd", cmd, "err", err)
		}
		return
	}

	g.jobQueue.Enqueue(job)
}
