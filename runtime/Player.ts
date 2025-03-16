class Player {
    private native: any
    Send(format: string, ...args: string[]) {
        this.native.Sendf(format, ...args)
    }
}