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
declare class Link {
}
declare class Player {
    private native;
    Send(format: string, ...args: string[]): void;
}
declare class Room {
}
