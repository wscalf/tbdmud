class Lockable extends Link {
    private locked: boolean = true
    private key: string = ""

    override Move(player: Player, to: Room): void {
        if (!this.locked)
            super.Move(player, to);
        else
            player.Send("The door is locked.");
    }

    Unlock(player: Player) {
        this.locked = false;
        player.Send("You unlock the door with a heavy clank.")
        player.Room.SendToAllExcept(player, "%s unlocks the door with a heavy clank.", player.Name);
    }

    Lock(player: Player) {
        this.locked = true;
        player.Send("You lock the door with a heavy clank.")
        player.Room.SendToAllExcept(player, "%s locks the door with a heavy clank.", player.Name);
    }
}