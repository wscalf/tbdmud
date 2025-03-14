package main

import (
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/net"
)

func main() {
	defaultPort := 4000      //Make parameter
	worldPath := "../sample" //Make parameter

	loader := game.NewLoader(worldPath)
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
	game := game.NewGame(commands, telnetListener, world, login, layouts)

	game.Run()
}
