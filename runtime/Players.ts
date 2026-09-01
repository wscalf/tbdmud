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

    private _on_player_joined(player: Player) {
        this.Connected.emit(player);
    }
    public Connected: Signal<[Player]> = new Signal();

    private _on_player_leaving(player: Player) {
        this.Disconnecting.emit(player);
    }
    public Disconnecting: Signal<[Player]> = new Signal();
}

declare const Players: _Players