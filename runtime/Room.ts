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

    SendToAll(pattern: string, ...args: string[]) {
        this.native.SendToAll(pattern, ...args);
    }

    SendToAllExcept(player: Player, pattern: string, ...args: string[]) {
        this.native.SendToAllExcept(player["native"], pattern, ...args);
    }

    FindPathTo(to: Room, limit: number): Link[] | null {
        let [path, found]: [any[], boolean] = this.native.FindPathTo(to.native, limit);
        if (!found) {
            return null;
        }

        let result: Link[] = []
        path.forEach(l => {
            result.push(extractJSObj(l))
        });

        return result;
    }
}