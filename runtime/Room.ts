class Room {
    private native: any;

    get Name(): string {
        return this.native.Name
    }

    set Name(value: string) {
        this.native.Name = value;
    }

    get Players(): Player[] {
        let players: Player[] = [];
        
        for (var player of this.native.GetPlayers()) {
            players.push(extractJSObj(player))
        }

        return players;
    }

    get Links(): Link[] {
        let links: Link[] = [];

        for (var link of this.native.GetLinks()) {
            links.push(extractJSObj(link))
        }

        return links
    }
}