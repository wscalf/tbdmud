class LockCommands {
    @Command("unlock", "Unlocks the given <exit>", [{name: "exit", type: "name", required: true}])
    static unlock(player: Player, link: string) {
        let l = LockCommands.findExit(link, player.Room);
        if (l == null) {
            player.Send("I don't see that here.")
        } else if (l instanceof Lockable) {
            player.Send("Trying to unlock");
            l.Unlock(player)
        } else {
            player.Send("You can't unlock that.")
        }
    }

    @Command("lock", "locks the given <exit>", [{name: "exit", type: "name", required: true}])
    static lock(player: Player, link: string) {
        let l = LockCommands.findExit(link, player.Room);
        if (l == null) {
            player.Send("I don't see that here.")
        } else if (l instanceof Lockable) {
            l.Lock(player)
        } else {
            player.Send("You can't lock that.")
        }
    }

    static findExit(name: string, room: Room): Link | null {
        for (var link of room.Links) {
            if (link.Command == name) {
                return link
            }
        }

        return null;
    }
}