package game

import (
	"errors"
	"log/slog"

	"github.com/wscalf/tbdmud/internal/text"
)

// Might not need these to be global since everything goes through Game anyway
var _scriptSystem ScriptSystem
var _world *World

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
	_scriptSystem = scriptSystem
	_world = world

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
			data, err := g.login.Process(client)
			if err != nil {
				slog.Error("error processing login", "err", err)
				client.Send(text.NewPrintfJob("An error occurred.")) //Consider looping back to the login process- or even looping the specific step
				return
			}
			if data.ID == "" {
				client.Disconnect()
				return
			}

			p := NewPlayer(data.ID, data.Name)
			p.Description = data.Desc

			g.players[p.ID] = p
			p.AttachClient(client)
			p.SetInputHandler(g.handleCommand)
			p.SetLayout(g.login.defaultPlayerLayout)

			g.jobQueue.Enqueue(JoinWorldJob{
				player: p,
				data:   data,
				game:   g,
			})

			p.Run()
			data, err = p.GetSaveData() //Probably shouldn't keep the data object around for the whole player session - consider some scoping
			if err != nil {
				slog.Error("error getting player save data", "err", err)
			} else {
				g.login.store.CreateOrUpdatePlayer(data)
			}

			p.Leave()
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
	data   *PlayerSaveData
	game   *Game
}

func (j JoinWorldJob) Run() {
	script, err := _scriptSystem.Wrap(j.player, j.game.defaultPlayerType)
	if err != nil {
		slog.Error("error wrapping character script", "err", err, "script", j.game.defaultPlayerType)
		j.player.client.Disconnect()
		return
	}
	for name, value := range j.data.Vars {
		script.Set(name, value)
	}
	j.player.AttachScript(script)

	if j.data.RoomID != "" {
		if room, ok := j.game.world.rooms[j.data.RoomID]; ok {
			j.player.Join(room)
		} else {
			j.player.Join(j.game.world.defaultRoom)
		}
	} else {
		j.player.Join(j.game.world.chargen)
	}
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
