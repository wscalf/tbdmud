package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/net"
	"github.com/wscalf/tbdmud/internal/scripting"
)

func main() {
	portValue := os.Getenv("TELNET_PORT")
	port, err := strconv.Atoi(portValue)
	if err != nil {
		slog.Error("Failed to parse TELNET_PORT as integer", "value", portValue, "err", err)
		return
	}
	worldPath := os.Getenv("WORLD")

	loader := game.NewLoader(worldPath)

	scriptSystem, err := initializeScripting(loader)
	if err != nil {
		slog.Error("Failed to initialize scripting subsystem. Exiting..", "err", err)
		return
	}

	world := game.NewWorld()
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

	rooms, err := loader.GetRooms(scriptSystem, meta.DefaultRoomType, meta.DefaultLinkType, meta.DefaultObjectType) //Get default type names from metadata
	if err != nil {
		slog.Error("Failed to load rooms. Exiting..", "err", err)
		return
	}

	world.SetRoomLayout(layouts["room"])
	world.InitializeRooms(rooms)
	world.SetSystemRooms(meta.ChargenRoom, meta.DefaultRoom)

	login := game.NewLogin(meta.Banner, layouts["player"])

	telnetListener := net.NewTelnetListener(port)
	commands := game.NewCommands()
	commands.RegisterBuiltins(layouts)

	scriptSystem.RegisterCommands(commands)

	game := game.NewGame(commands, telnetListener, world, login, layouts, scriptSystem, meta.DefaultPlayerType)

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
