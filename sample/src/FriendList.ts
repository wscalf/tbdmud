class FriendData {
    constructor(name: string, note: string) {
        this.Name = name;
        this.Note = note;
    }

    Name: string
    Note: string
}

class FriendList extends Dict<FriendData> {
    static accessor: (p: Player) => FriendList;

    @Command("add-friend", "Adds another character to your friend list", [{name: "name", type: "name", required: true}, {name: "note", type: "freetext", required: false}])
    static add_friend(player: Player, name: string, note: string) {
        let other = Players.FindByName(name); //Note: this will currently only work for online players
        if (other != null) {
            let list = FriendList.accessor(player);
            list.Add(other, note);
            player.Send("%s has been added to your friend list.", name);
        } else {
            player.Send("%s was not found.", name);
        }
    }

    @Command("remove-friend", "Adds another character to your friend list", [{name: "name", type: "name", required: true}])
    static remove_friend(player: Player, name: string) {
        let other = Players.FindByName(name); //Note: this will currently only work for online players
        if (other != null) {
            let list = FriendList.accessor(player);
            list.Remove(other);
            player.Send("%s was removed from your friend list.", name);
        } else {
            player.Send("%s was not found.", name);
        }
    }

    @Command("list-friends", "Lists online friends", [])
    static list_friends(player: Player) {
        let list = FriendList.accessor(player);
        player.Send("------------- Online Friends -------------")
        player.Send("==========================================")
        list.forEach((id, data) => {
            let other = Players.FindById(id)
            if (other == null) {
                return; //Skip
            }

            if (data.Note) {
                player.Send("%s:\t\t\t%s", data.Name, data.Note)
            } else {
                player.Send("%s", data.Name)
            }
        })
    }

    public Add(p: Player, note: string) {
        this.set(p.ID, new FriendData(p.Name, note));
    }

    public Remove(p: Player) {
        this.remove(p.ID);
    }

    static notify_connected_players(changed: Player, message: string) {
        Log.Debug(`Notifyng of change for ${changed.Name}`)
        Players.All().forEach(p => {
            if (p == changed) {
                return; //Skip
            }
    
            let list = FriendList.accessor(p);
            if (list.has(changed.ID)) {
                p.Send(message);
            }
        });
    }
}

function initialize_friends_list(accessor: (p: Player) => FriendList, joinTemplate: (p: Player) => string, leavingTemplate: (p: Player) => string) {
    FriendList.accessor = accessor;

    Players.Connected.connect((p: Player) => {
        let message = joinTemplate(p);
        FriendList.notify_connected_players(p, message);
    })

    Players.Disconnecting.connect((p: Player) => {
        let message = leavingTemplate(p);
        FriendList.notify_connected_players(p, message);
    })
}

