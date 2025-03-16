class Link {
    private native: any;

    get Name(): string {
        return this.native.Name;
    }
    set Name(value: string) {
        this.native.Name = value;
    }

    get Command(): string {
        return this.native.Command;
    }

    Move(player: Player, to: Room) {
        this.native.Move(player["native"], to["native"]);
    }
}