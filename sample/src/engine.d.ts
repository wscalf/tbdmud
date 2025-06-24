declare type Parameter = {
    name: string;
    required: boolean;
    type: "name" | "freetext";
};
declare type CommandRef = {
    name: string;
    desc: string;
    params: Parameter[];
    handler: Function;
};
declare const commands: CommandRef[];
declare function Command(name: string, desc: string, params: Parameter[]): (target: any, key: string, descriptor: PropertyDescriptor) => void;
declare class ChatMessage {
    Actor: string;
    Body: string;
}
declare class _GenAI {
    private native;
    Generate(systemPrompt: string, prompt: string): string;
    Chat(systemPrompt: string, history: ChatMessage[], actor: string, pose: string): ChatMessage[];
}
declare const GenAI: _GenAI;
declare function extractJSObj(native: any): any;
declare class Link {
    private native;
    get Name(): string;
    set Name(value: string);
    get Command(): string;
    Peek(): Room;
    Move(player: Player, to: Room): void;
}
declare class MUDObject {
    private native;
    get ID(): string;
    get Name(): string;
    set Name(value: string);
    get Desc(): string;
    set Desc(value: string);
}
declare class Map<T> {
    private data;
    has(key: string): boolean;
    get(key: string): T;
    set(key: string, value: T): void;
    remove(key: string): void;
    forEach(callbackFn: (key: string, value: T) => void): void;
}
declare class Player {
    private native;
    get Name(): string;
    set Name(value: string);
    get Desc(): string;
    set Desc(value: string);
    get Room(): Room;
    Send(format: string, ...args: string[]): void;
    get Items(): MUDObject[];
}
declare class _Players {
    private native;
    FindById(id: string): Player | null;
    FindByName(name: string): Player | null;
    All(): Player[];
}
declare const Players: _Players;
declare let persistedPropertiesByType: Map<Array<string>>;
declare function persist(): (proto: any, member: string) => void;
declare function getPersistedProperties(typeName: string): Array<string>;
declare class Room {
    private native;
    get ID(): string;
    get Name(): string;
    set Name(value: string);
    get Desc(): string;
    get Players(): Player[];
    get Links(): Link[];
    FindItem(name: string): MUDObject;
    SendToAll(pattern: string, ...args: string[]): void;
    SendToAllExcept(player: Player, pattern: string, ...args: string[]): void;
    FindPathTo(to: Room, limit: number): Link[] | null;
}
declare class _World {
    private native;
    FindRoom(id: string): Room | null;
}
declare const World: _World;
