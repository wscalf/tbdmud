class _Players {
    private native: any

    public FindById(id: string): Player | null {
        let p: Player | null = this.native.FindById(id);
        if (p == null) return null;

        return extractJSObj(p);
    }

    public FindByName(name: string): Player | null {
        let p: Player | null = this.native.FindByName(name);
        if (p == null) return null;

        return extractJSObj(p);
    }

    public All(): Player[] {
        let nativePlayers: any[] = this.native.All();
        let players: Player[] = [];

        nativePlayers.forEach(p =>
            players.push(extractJSObj(p))
        )

        return players;
    }
}

declare const Players: _Players