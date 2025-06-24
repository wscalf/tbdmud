class Map<T> {
    private data: {[key: string]: T} = {};
    public has(key: string): boolean {
        return this.data.hasOwnProperty(key);
    }

    public get(key: string): T {
        return this.data[key];
    }

    public set(key: string, value: T) {
        this.data[key] = value;
    }

    public remove(key: string) {
        delete this.data[key];
    }

    public forEach(callbackFn: (key: string, value: T) => void) {
        Object.keys(this.data).forEach(key => {
            let value = this.data[key];

            callbackFn(key, value);
        })
    }
}