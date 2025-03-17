package game

type Object struct {
	ID          string
	Name        string
	Description string
	script      ScriptObject
}

func NewObject(name, description string) *Object {
	return &Object{
		Name:        name,
		Description: description,
	}
}

func (o *Object) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"name": o.Name,
		"desc": o.Description,
	}
}

func (o *Object) GetScript() ScriptObject {
	return o.script
}

func (o *Object) AttachScript(script ScriptObject) {
	o.script = script
}
