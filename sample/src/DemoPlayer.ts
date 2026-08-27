class DemoPlayer extends Player {
    @persist()
    public Count: number = 0;

    constructor() {
        super()

        Players.Connected.connect(this.handlePlayerJoined);
        Players.Disconnecting.connect(this.handlePlayerDisconnecting);
    }
 
    private handlePlayerJoined = (player: Player) => {
        player = player as DemoPlayer;
        if (player instanceof DemoPlayer) {
            if (player != this) {
                this.Send("%s joined...", player.Name);
            }
        }
    }

    private handlePlayerDisconnecting = (player: Player) => {
        if (this == player) {
            Players.Connected.disconnect(this.handlePlayerJoined);
            Players.Disconnecting.disconnect(this.handlePlayerDisconnecting);
        }
    }
}