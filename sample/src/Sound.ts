class _Sound {
    EmitLoud(from: Room, limit: number, pose: string) {
        const systemPrompt: string = `The world is: ${WorldDesc}. You are a storyteller in a roleplaying game that's already in progress. You need to describe a single event to your players taking into account the rooms between where it occurs and where they are. Keep your descriptions brief and focused on the event, and do not mention anyone or anything not explicitly described to you.`
        //This could probably be simplified and made more scalable with a visitor pattern
        //For hackathon purposes, sticking with what we have

        let players: Player[] = Players.All();
        let occupiedRooms = new Map<Room>;

        players.forEach(p => {
            let room = p.Room
            if (!occupiedRooms.has(room.ID)) {
                occupiedRooms.set(room.ID, room);
            }
        })

        from.SendToAll(pose);

        occupiedRooms.forEach((_, room) => {
            if (from.ID != room.ID) {
                let path = from.FindPathTo(room, limit);
                if (path == null || path.length == 0) {
                    return; //No path found, or its empty.
                }

                let prompt = `The event takes place in ${from.Name}, which is: ${from.Desc}\n\nThe event is: ${pose}\n\n The shortest path between that and your players is:\n`
                path.forEach(link => {
                    let next = link.Peek();
                    if (next.ID != room.ID)
                        prompt += `${next.Name}, which is: ${next.Desc}\n\n`;
                })

                prompt += `Finally, your players are in ${room.Name}, which is: ${room.Desc}\n\nDescribe the event to them.`

                let emit = GenAI.Generate(systemPrompt, prompt);

                room.SendToAll(emit);
            }
        })
    }
}

const Sound: _Sound = new _Sound();