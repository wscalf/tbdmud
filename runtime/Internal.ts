function extractJSObj(native: any): any {
    let scriptObj: any = native.GetScript() //All native objects impl GetScript()
    return scriptObj.Obj //The JS object is stored in ScriptObject at "Obj" 
}