package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/net"
	"github.com/wscalf/tbdmud/internal/scripting"
	"github.com/wscalf/tbdmud/internal/storage"
)

func main() {
	portValue := os.Getenv("TELNET_PORT")
	port, err := strconv.Atoi(portValue)
	if err != nil {
		slog.Error("Failed to parse TELNET_PORT as integer", "value", portValue, "err", err)
		return
	}
	worldPath := os.Getenv("WORLD")
	store, err := initializeStorage(worldPath)
	if err != nil {
		slog.Error("Error initializing storage provider.", "err", err)
		return
	}

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

	login := game.NewLogin(meta.Banner, layouts["player"], store)

	telnetListener := net.NewTelnetListener(port)
	commands := game.NewCommands()
	commands.RegisterBuiltins(layouts)

	scriptSystem.RegisterCommands(commands)

	game := game.NewGame(commands, telnetListener, world, login, layouts, scriptSystem, meta.DefaultPlayerType)

	game.Run()
}

func initializeStorage(worldPath string) (game.Storage, error) {
	store := storage.NewBoltDBStore()
	err := store.Initialize(worldPath)
	if err != nil {
		return nil, err
	}

	go store.Process()
	return store, nil
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

	err = system.Initialize()
	if err != nil {
		return nil, fmt.Errorf("error initialize script system: %w", err)
	}

	return system, nil
}
