class Signal<TArgs extends any[]> {
    private handlers: Array<(...args: TArgs) => void> = [];
    
    connect(handler: (...args: TArgs) => void): void {
        this.handlers.push(handler);
    }

    disconnect(handler: (...args: TArgs) => void): void {
        var index = this.handlers.indexOf(handler);

        if (index !== -1) {
            this.handlers.splice(index, 1);
        } else {
            Log.Debug("handler not found");
        }
    }

    emit(...args: TArgs): void {
        for (var i = 0; i < this.handlers.length; i++) {
            this.handlers[i].apply(null, args);
        }
    }
}