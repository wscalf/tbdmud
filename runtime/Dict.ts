class Dict<T> {
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

    public forEach(callback: (key: string, value: T) => void) {
        let keys = Object.getOwnPropertyNames(this.data);
        let sorted = keys.sort();

        sorted.forEach((k) => {
            let v = this.get(k);

            callback(k, v);
        })
    }
}