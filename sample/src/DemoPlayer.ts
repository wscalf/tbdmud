/// <reference path="FriendList.ts" />

class DemoPlayer extends Player {
    @persist()
    public Count: number = 0;

    public FriendList: FriendList = new FriendList();

}

initialize_friends_list((p) => {
    if (p instanceof DemoPlayer) {
        return p.FriendList;
    } else {
        throw new Error("Object not of type DemoPlayer");
    }
}, (p) => {
    return `${p.Name} has connected.`
}, (p) => {
    return `${p.Name} has disconnected.`
})