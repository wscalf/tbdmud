class Player {
    private native: any

    get ID(): string {
        return this.native.ID;
    }

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
        Log.Debug("Sending a message to: " + this.Name);
        this.native.Sendf(format, ...args)
    }

    get Items(): MUDObject[] {
        let items: MUDObject[] = []
        for (var item of this.native.GetItems()) {
            items.push(extractJSObj(item))
        }
        return items
    }
}