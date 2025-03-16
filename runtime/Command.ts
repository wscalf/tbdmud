type Parameter = {name: string, required: boolean, type: "name" | "freetext"}
type CommandRef = {name: string, desc: string, params: Parameter[], handler: Function}

const commands: CommandRef[] = [];

function Command(name: string, desc: string, params: Parameter[]) {
    return function (target: any, key: string, descriptor: PropertyDescriptor) {
        commands.push({name, desc, params, handler: descriptor.value})
    }
}