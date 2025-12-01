class PlayerSaveData {
    private native: any;

    get ID(): string {
        return this.native.ID;
    }
    get RoomID(): string {
        return this.native.RoomID;
    }
    get Name(): string {
        return this.native.Name;
    }
    get Desc(): string {
        return this.native.Desc;
    }
    get Vars(): Map<any> {
        return this.native.Vars;
    }
}