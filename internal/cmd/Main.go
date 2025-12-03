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
	"github.com/wscalf/tbdmud/web"
)

func main() {
	telnetPortValue := os.Getenv("TELNET_PORT")
	telnetPort, err := strconv.Atoi(telnetPortValue)
	if err != nil {
		slog.Error("Failed to parse TELNET_PORT as integer", "value", telnetPortValue, "err", err)
		return
	}

	httpPortValue := os.Getenv("HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortValue)
	if err != nil {
		slog.Error("Failed to parse HTTP_PORT as integer", "value", httpPortValue, "err", err)
		return
	}

	worldPath := os.Getenv("WORLD")
	store, err := initializeStorage(worldPath)
	if err != nil {
		slog.Error("Error initializing storage provider.", "err", err)
		return
	}

	loader := game.NewLoader(worldPath)
	world := game.NewWorld()

	players := game.NewPlayers()

	scriptSystem, err := initializeScripting(loader, world, players)
	if err != nil {
		slog.Error("Failed to initialize scripting subsystem. Exiting..", "err", err)
		return
	}

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

	listeners := net.NewAggregateClientListener()
	telnetListener := net.NewTelnetListener(telnetPort)
	listeners.AddListener(telnetListener)
	web := web.NewWeb()
	listeners.AddListener(web.Listener())

	commands := game.NewCommands()
	commands.RegisterBuiltins(layouts)

	scriptSystem.RegisterCommands(commands)

	game := game.NewGame(commands, listeners, players, world, login, layouts, scriptSystem, meta.DefaultPlayerType)

	web.ServeHttp(httpPort)
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

func initializeScripting(loader *game.Loader, world *game.World, players *game.Players) (game.ScriptSystem, error) {
	system := scripting.NewGojaScriptSystem()

	err := system.RunBootstrapCode()
	if err != nil {
		return nil, fmt.Errorf("error executing engine.js: %w", err)
	}

	moduleCode, err := loader.ReadModuleTextFile("module.js")
	if err != nil {
		return nil, fmt.Errorf("error reading module.js: %w", err)
	}

	err = system.AddGlobal("World", "_World", world)
	if err != nil {
		return nil, fmt.Errorf("error binding World global: %w", err)
	}

	err = system.AddGlobal("Players", "_Players", players)
	if err != nil {
		return nil, fmt.Errorf("error binding Players global: %w", err)
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
