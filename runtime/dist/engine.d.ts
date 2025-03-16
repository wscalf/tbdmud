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
declare function extractJSObj(native: any): any;
declare class Link {
    private native;
    get Name(): string;
    set Name(value: string);
    Move(player: Player, to: Room): void;
}
declare class Player {
    private native;
    get Name(): string;
    set Name(value: string);
    get Room(): Room;
    Send(format: string, ...args: string[]): void;
}
declare class Room {
    private native;
    get Name(): string;
    set Name(value: string);
    get Players(): Player[];
    get Links(): Link[];
}
