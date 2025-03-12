package main

import (
	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/game/commands"
	"github.com/wscalf/tbdmud/internal/net"
)

func main() {
	defaultPort := 4000

	telnetListener := net.NewTelnetListener(defaultPort)
	commands := commands.NewCommands()
	commands.RegisterBuiltins()
	game := game.NewGame(commands, telnetListener)

	game.Run()
}
