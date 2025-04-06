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

func (o *Object) GetSaveData() (ObjectSaveData, error) {
	data := ObjectSaveData{
		ID:   o.ID,
		Name: o.Name,
		Desc: o.Description,
	}

	if o.script != nil {
		vars, err := o.script.GetSaveProperties()
		if err != nil {
			return data, err
		}

		data.TypeName = o.script.Type()
		data.Vars = vars
	}

	return data, nil
}

type ObjectSaveData struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Desc     string         `json:"desc"`
	TypeName string         `json:"type"`
	Vars     map[string]any `json:"vars"`
}

func ObjectFromSaveData(data map[string]any) Object {
	//Need to construct a script object from the stored type - singletonize ScriptSystem? Pass it everywhere?
	obj := Object{}

	if id, ok := data["id"]; ok {
		obj.ID = id.(string)
	}

	if name, ok := data["name"]; ok {
		obj.Name = name.(string)
	}

	if desc, ok := data["desc"]; ok {
		obj.Description = desc.(string)
	}

	return obj
}
