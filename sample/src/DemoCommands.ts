class DemoCommands {
    @Command("increment", "Increments a value on the player", [])
    static increment(player: DemoPlayer) {
        player.Count++;
    }

    @Command("check", "Checks the value on the player", [])
    static check(player: DemoPlayer) {
        player.Send("The current value is: %s", player.Count.toString());
    }
}