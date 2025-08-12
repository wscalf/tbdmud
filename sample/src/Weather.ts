const weatherSystemMessage = `The world is: ${WorldDesc}. You are the storyteller in a roleplaying game.`
class _Weather {
    private weather = "The skies are clear outside.";
    public SetWeather(description: string) {
        let players: Player[] = Players.All();
        let occupiedRooms = new Map<Room>;

        players.forEach(p => {
            let room = p.Room
            if (!occupiedRooms.has(room.ID)) {
                occupiedRooms.set(room.ID, room);
            }
        })

        occupiedRooms.forEach((_, room) => {
            let emit = GenAI.Generate(weatherSystemMessage, `Your players are in the room: ${room.Name}, which is: ${room.Desc}\n\nThe weather is now: ${description}\nDescribe the weather to your players.`);
            room.SendToAll(emit);
        })
    }
}

const Weather: _Weather = new _Weather();

class WeatherCommands {
    @Command("rain", "Causes it to start raining", [])
    static startRain() {
        Weather.SetWeather("It starts raining outside.");
    }

    @Command("rain", "Causes it to rain heavily", [])
    static heavyRain() {
        Weather.SetWeather("It's raining heavily outside.");
    }

    @Command("rain", "Causes it to stop raining", [])
    static clearSkies() {
        Weather.SetWeather("The skies are clear outside.");
    }
}