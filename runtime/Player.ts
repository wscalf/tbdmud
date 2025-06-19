class Player {
    private native: any

    get Name(): string {
        return this.native.Name;
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

    get Room(): Room {
        return extractJSObj(this.native.GetRoom());
    }

    Send(format: string, ...args: string[]) {
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