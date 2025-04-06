/// <reference path="Map.ts" />

let persistedPropertiesByType: Map<Array<string>> = new Map<Array<string>>();

function persist()
{
    return (proto: any, member: string) =>
    {
        let typeName = proto.constructor.name;
        if (persistedPropertiesByType.has(typeName)) {
            persistedPropertiesByType.get(typeName).push(member);
        } else {
            persistedPropertiesByType.set(typeName, new Array<string>(member));
        }
    }
}

function getPersistedProperties(typeName: string): Array<string> {
    if (persistedPropertiesByType.has(typeName)) {
        return persistedPropertiesByType.get(typeName);
    } else {
        return new Array<string>();
    }
}