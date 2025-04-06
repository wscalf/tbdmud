package scripting

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

var ErrUnrecognizedProperty error = errors.New("unrecognized property")

type GojaScriptObject struct {
	Obj        *goja.Object
	system     *GojaScriptSystem
	scriptType string
	//propertyNames map[string]bool
}

func newGojaScriptObject(obj *goja.Object, system *GojaScriptSystem, typeName string) *GojaScriptObject {
	o := &GojaScriptObject{
		Obj:        obj,
		system:     system,
		scriptType: typeName,
		//propertyNames: map[string]bool{},
	}

	//Checking property names is unreliable - they don't show up until assigned - possible ECMAScript 5 issue
	// for _, name := range obj.GetOwnPropertyNames() {
	// 	o.propertyNames[name] = true
	// }

	return o
}

func (o *GojaScriptObject) Get(prop string) (interface{}, error) {
	//Checking property names is unreliable - they don't show up until assigned - possible ECMAScript 5 issue
	// if !o.propertyNames[prop] {
	// 	return nil, fmt.Errorf("unable to access property %s: %w", prop, ErrUnrecognizedProperty)
	// }

	v := o.Obj.Get(prop)
	return o.system.exportValue(v), nil
}

func (o *GojaScriptObject) Set(prop string, value interface{}) error {
	//Checking property names is unreliable - they don't show up until assigned - possible ECMAScript 5 issue
	// if !o.propertyNames[prop] {
	// 	return fmt.Errorf("unable to access property %s: %w", prop, ErrUnrecognizedProperty)
	// }

	err := o.Obj.Set(prop, o.system.importValue(value))
	if err != nil {
		return fmt.Errorf("error setting property %s: %w", prop, err)
	}

	return nil
}

func (o *GojaScriptObject) Call(name string, args ...interface{}) (interface{}, error) {
	value := o.Obj.Get(name)
	if goja.IsUndefined(value) {
		return nil, fmt.Errorf("error calling method %s on object of type %s: method not found", name, o.Obj.ClassName())
	}

	method, ok := goja.AssertFunction(value)
	if !ok {
		return nil, fmt.Errorf("error calling method %s on object of type %s: not a method", name, o.Obj.ClassName())
	}

	importedArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		importedArgs[i] = o.system.importValue(arg)
	}

	v, err := method(o.Obj, importedArgs...)
	if err != nil {
		return nil, fmt.Errorf("error calling method %s with args %v: %w", name, args, err)
	}

	return o.system.exportValue(v), nil
}

func (o *GojaScriptObject) Type() string {
	return o.scriptType
}

func (o *GojaScriptObject) GetDescribeProperties() (map[string]any, error) {
	return map[string]any{}, nil
}

func (o *GojaScriptObject) GetSaveProperties() (map[string]any, error) {
	fields, err := o.system.getPersistedFieldsForType(o.scriptType)
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	for _, field := range fields {
		v, err := o.Get(field)
		if err != nil {
			return nil, err
		}

		data[field] = v
	}

	return data, nil
}
