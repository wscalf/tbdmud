class MUDObject {
    private native: any

    get Name(): string {
        return this.native.Name
    }
    set Name(value: string) {
        this.native.Name = value;
    }

    get Desc(): string {
        return this.native.Description
    }
    set Desc(value: string) {
        this.native.Description = value;
    }
}