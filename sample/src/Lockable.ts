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
        let usedKey = this.findMatchingKey(player);
        if (usedKey == null) {
            player.Send("You don't have the key.")
            return
        }

        this.locked = false;
        player.Send("You unlock the door with the %s.", usedKey.Name)
        player.Room.SendToAllExcept(player, "%s unlocks the door with a the %s.", player.Name, usedKey.Name);
    }

    Lock(player: Player) {
        let usedKey = this.findMatchingKey(player);
        if (usedKey == null) {
            player.Send("You don't have the key.")
            return
        }

        this.locked = true;
        player.Send("You lock the door with the %s.", usedKey.Name)
        player.Room.SendToAllExcept(player, "%s locks the door with a the %s.", player.Name, usedKey.Name);
    }

    private findMatchingKey(player: Player): Key | null {
        for (var item of player.Items) {
            if (item instanceof Key && item.Matches(this.key)) {
                return item
            }
        }

        return null
    }
}