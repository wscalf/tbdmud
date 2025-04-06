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
}