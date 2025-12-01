class _System {
    private native: any;

    public Wait(seconds: number): Promise<void> {
        let promise: Promise<any> = this.native.Wait(seconds);
        return promise;
    }
}

declare const System: _System;