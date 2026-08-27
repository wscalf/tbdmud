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

    public async FindByNameIncludingOffline(name: string): Promise<PlayerSaveData> {
        let p: any = await this.native.FindByNameIncludingOffline(name);
        return new PlayerSaveData();
    }
}

declare const Players: _Players