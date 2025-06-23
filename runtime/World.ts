class _World {
    private native: any;

    public FindRoom(id: string): Room | null {
        let found: Room = this.native.FindRoom(id);
        if (!found) return null;

        return extractJSObj(found);
    }
}

declare const World: _World;