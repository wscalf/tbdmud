class TestCommands {
    @Command("echo", "Sends <message> to the player using the command.", [{name: "message", type: "freetext", required: true}])
    static echo(player: Player, message: string): void {
        player.Send(message)
    }
}
