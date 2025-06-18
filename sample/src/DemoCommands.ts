class DemoCommands {
    @Command("increment", "Increments a value on the player", [])
    static increment(player: DemoPlayer) {
        player.Count++;
    }

    @Command("check", "Checks the value on the player", [])
    static check(player: DemoPlayer) {
        player.Send("The current value is: %s", player.Count.toString());
    }

    @Command("pathTo", "Finds the shortest path to <roomId>", [{name: "roomId", type: "name", required: true}, {name: "limit", type: "name", required: true}])
    static pathTo(player: DemoPlayer, roomId: string, limit: number) {
        let other = World.FindRoom(roomId)
        if (!other) {
            player.Send("There is no room with the id %s.", roomId);
            return;
        }

        let path = player.Room.FindPathTo(other, limit);
        if (path == null) {
            player.Send("No path found to room %s within %s moves.", roomId, limit.toString());
        } else if (path.length == 0) {
            player.Send("The path is zero moves - you're already there!");
        } else {
            path.forEach(link => {
                let cmd: string = link.Command;
                let destination: Room = link.Peek();

                player.Send("Go %s to %s", cmd, destination.Name);
            });
        }
    }
    @Command("test-ai", "Prompts the integrated generative AI", [{name: "prompt", type: "freetext", required: true}])
    static test_ai(player: DemoPlayer, prompt: string) {
        let result: string = GenAI.Generate("You are AnyVAC, an advanced machine intelligence that assists humans by responding to prompts, inspired by the MI from Isaac Asimov's The Last Question.", prompt)

        player.Send(result)
    }
}