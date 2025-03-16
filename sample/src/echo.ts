class TestCommands {
    @Command("echo", "Sends <message> to the player using the command.", [{name: "message", type: "freetext", required: true}])
    static echo(player: Player, message: string): void {
        player.Send(message)
    }
    @Command("whoami", "Outputs the current player's name.", [])
    static whoami(player: Player) {
        player.Send(player.Name);
    }
    @Command("whereami", "Outputs the current player's room.", [])
    static whereami(player: Player) {
        player.Send(player.Room.Name);
    }
    @Command("who", "Outputs the players in the room.", [])
    static who(player: Player) {
        for (var p of player.Room.Players) {
            player.Send(p.Name);
        }
    }
    @Command("exits", "Lists exits from the current room.", [])
    static exits(player: Player) {
        for (var l of player.Room.Links) {
            player.Send(l.Name)
        }
    }
}
