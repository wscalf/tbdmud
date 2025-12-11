class DemoCommands {
    @Command("increment", "Increments a value on the player", [])
    static increment(player: DemoPlayer) {
        player.Count++;
    }

    @Command("check", "Checks the value on the player", [])
    static check(player: DemoPlayer) {
        player.Send("The current value is: %s", player.Count.toString());
    }

    @Command("check-other", "Checks the value on the other player", [{name: "other", type: "name", required: true}])
    static async checkOther(player: DemoPlayer, other: string) {
        player.Send("Checking other player...");
        let data = await Players.FindByNameIncludingOffline(other);
        player.Send("Back from await.");
        player.Send("The other player's name is: %s", data.Name);
        player.Send("The other player's value is: %s", data.Vars.get("Count"));
    }

    @Command("long-think", "Prints the message after 1 second.", [{name: "message", type: "freetext", required: true}])
    static async longThink(player: DemoPlayer, message: string) {
        await System.Wait(1);
        player.Send(message);
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
}