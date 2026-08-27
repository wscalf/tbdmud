class _Log {
    native: any;

    public Debug(format: string, ...args: string[]) {
        this.native.Debug(format, args);
    }
}

declare const Log: _Log;