package main

import (
	"fmt"
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/net"
	"github.com/wscalf/tbdmud/internal/scripting"
)

func main() {
	defaultPort := 4000     //Make parameter
	worldPath := "./sample" //Make parameter

	loader := game.NewLoader(worldPath)

	scriptSystem, err := initializeScripting(loader)
	if err != nil {
		slog.Error("Failed to initialize scripting subsystem. Exiting..", "err", err)
		return
	}

	world := game.NewWorld(scriptSystem, "Room")
	meta, err := loader.GetMeta()
	if err != nil {
		slog.Error("Failed to load module metadata. Exiting..", "err", err)
		return
	}

	layouts, err := loader.GetLayouts()
	if err != nil {
		slog.Error("Failed to load layouts. Exiting..", "err", err)
		return
	}

	rooms, err := loader.GetRooms()
	if err != nil {
		slog.Error("Failed to load rooms. Exiting..", "err", err)
		return
	}

	world.SetRoomLayout(layouts["room"])
	world.InitializeRooms(rooms)
	world.SetSystemRooms(meta.ChargenRoom, meta.DefaultRoom)

	login := game.NewLogin(meta.Banner, layouts["player"])

	telnetListener := net.NewTelnetListener(defaultPort)
	commands := game.NewCommands()
	commands.RegisterBuiltins()

	scriptSystem.RegisterCommands(commands)

	game := game.NewGame(commands, telnetListener, world, login, layouts, scriptSystem, "Player")

	game.Run()
}

func initializeScripting(loader *game.Loader) (game.ScriptSystem, error) {
	system := scripting.NewGojaScriptSystem()

	err := system.RunBootstrapCode()
	if err != nil {
		return nil, fmt.Errorf("error executing engine.js: %w", err)
	}

	moduleCode, err := loader.ReadModuleTextFile("module.js")
	if err != nil {
		return nil, fmt.Errorf("error reading module.js: %w", err)
	}

	err = system.Run(moduleCode)
	if err != nil {
		return nil, fmt.Errorf("error executing module.js: %w", err)
	}

	system.Initialize()

	return system, nil
}
