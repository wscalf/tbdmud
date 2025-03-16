package game

type Object struct {
	ID          string
	Name        string
	Description string
	script      ScriptObject
}

func (o *Object) GetScript() ScriptObject {
	return o.script
}

func (o *Object) AttachScript(script ScriptObject) {
	o.script = script
}
