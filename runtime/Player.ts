class Player {
    private native: any

    get Name(): string {
        return this.native.Name;
    }
    set Name(value: string) {
        this.native.Name = value;
    }

    get Room(): Room {
        return extractJSObj(this.native.GetRoom());
    }

    Send(format: string, ...args: string[]) {
        this.native.Sendf(format, ...args)
    }
}