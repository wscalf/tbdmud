package scripting

import (
	"errors"
	"fmt"

	"github.com/dop251/goja"
)

var ErrUnrecognizedProperty error = errors.New("unrecognized property")

type GojaScriptObject struct {
	obj    *goja.Object
	system *GojaScriptSystem
	//propertyNames map[string]bool
}

func newGojaScriptObject(obj *goja.Object, system *GojaScriptSystem) *GojaScriptObject {
	o := &GojaScriptObject{
		obj:    obj,
		system: system,
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

	v := o.obj.Get(prop)
	return o.system.exportValue(v), nil
}

func (o *GojaScriptObject) Set(prop string, value interface{}) error {
	//Checking property names is unreliable - they don't show up until assigned - possible ECMAScript 5 issue
	// if !o.propertyNames[prop] {
	// 	return fmt.Errorf("unable to access property %s: %w", prop, ErrUnrecognizedProperty)
	// }

	err := o.obj.Set(prop, o.system.importValue(value))
	if err != nil {
		return fmt.Errorf("error setting property %s: %w", prop, err)
	}

	return nil
}

func (o *GojaScriptObject) Call(name string, args ...interface{}) (interface{}, error) {
	value := o.obj.Get(name)
	if goja.IsUndefined(value) {
		return nil, fmt.Errorf("error calling method %s on object of type %s: method not found", name, o.obj.ClassName())
	}

	method, ok := goja.AssertFunction(value)
	if !ok {
		return nil, fmt.Errorf("error calling method %s on object of type %s: not a method", name, o.obj.ClassName())
	}

	importedArgs := make([]goja.Value, len(args))
	for i, arg := range args {
		importedArgs[i] = o.system.importValue(arg)
	}

	v, err := method(o.obj, importedArgs...)
	if err != nil {
		return nil, fmt.Errorf("error calling method %s with args %v: %w", name, args, err)
	}

	return o.system.exportValue(v), nil
}
