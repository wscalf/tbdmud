const WorldDesc = "Most people live in a system of underground tunnels that were built so long ago that no one remembers what they were originally for. Going to the surface is risky but can be done, usually for short raids for supplies or treasure. No one really remembers why it's unsafe, though there are a range of myths depending on what appeals to any particular person or group. Some tunnels are well explored and mapped, and are generally considered safe, while others are not. Unexplored tunnels are unpredictable and may be unstable, confusing, full of hazardous chemicals, or even subject to collapse. People live in small, dispersed groups, each with their own distinct cultures, trying to survive who may at times come into competition or conflict, but also form bonds and friendly rivalries.";

class NPC extends MUDObject {
    private history: ChatMessage[] = [];

    Address(room: Room, pose: string) {
        let results = GenAI.Chat(this.createSystemPrompt(room), this.history, "", pose);
        let response = results[results.length - 1];
        this.history.push(...results)

        room.SendToAll(response.Body);
    }

    private createSystemPrompt(room: Room): string {
        let prompt: string = `
        The world is: ${WorldDesc}
        You are in: ${room.Desc}
        You are: ${this.Desc}
        `

        room.Players.forEach(p => {
            prompt += `${p.Name} is ${p.Desc}`
        });

        prompt += "You will roleplay as the given character. You will always stay in character and never break character or speak for other characters. You will continue conversations and situations in a natural fashion."
        return prompt
    }
}