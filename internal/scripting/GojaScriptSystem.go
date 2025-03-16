package scripting

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/dop251/goja"
	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

var ErrUnrecognizedType = errors.New("unrecognized type")

//go:embed engine.js
var bootstrapCode string

type GojaScriptSystem struct {
	vm    *goja.Runtime
	types map[string]goja.Value
}

func NewGojaScriptSystem() *GojaScriptSystem {
	return &GojaScriptSystem{
		vm:    goja.New(),
		types: map[string]goja.Value{},
	}
}

func (s *GojaScriptSystem) RunBootstrapCode() error {
	return s.Run(bootstrapCode)
}

func (s *GojaScriptSystem) Run(src string) error {
	_, err := s.vm.RunString(src)

	return err
}

func (s *GojaScriptSystem) Wrap(native interface{}, typeName string) (game.ScriptObject, error) {
	obj, err := s.createObject(typeName)
	if err != nil {
		return nil, err
	}
	obj.obj.Set("native", native) //Bypass the setter for native pointer so it doesn't try to unwrap it
	return obj, nil
}

func (s *GojaScriptSystem) Initialize() {
	//Find all the types
	for _, name := range s.vm.GlobalObject().GetOwnPropertyNames() {
		candidate := s.vm.Get(name)
		if _, ok := goja.AssertConstructor(candidate); ok {
			s.types[name] = candidate
		}
	}
}

func (s *GojaScriptSystem) RegisterCommands(commands *game.Commands) {
	refs := s.vm.Get("commands").(*goja.Object)
	for _, key := range refs.Keys() {
		ref := refs.Get(key).(*goja.Object)

		name := ref.Get("name").String()
		description := ref.Get("desc").String()
		handler, _ := goja.AssertFunction(ref.Get("handler"))
		params := ref.Get("params").(*goja.Object)

		cmd := NewGojaScriptCommand(name, description, s.convertParameters(params), s, handler)
		commands.Register(name, cmd)
	}
}

func (s *GojaScriptSystem) convertParameters(paramsObject *goja.Object) []parameters.Parameter {
	keys := paramsObject.Keys()
	params := make([]parameters.Parameter, 0, len(keys))
	for _, key := range keys {
		obj := paramsObject.Get(key).(*goja.Object)

		name := obj.Get("name").String()
		required := obj.Get("required").ToBoolean()
		paramType := obj.Get("type").String()

		switch paramType {
		case "name":
			params = append(params, parameters.NewName(name, required))
		case "freetext":
			params = append(params, parameters.NewFreeText(name))
		}
	}

	return params
}

func (s *GojaScriptSystem) createObject(typeName string) (*GojaScriptObject, error) {
	constructor, ok := s.types[typeName]
	if !ok {
		return nil, fmt.Errorf("unable to instantiate object of type %s: %w", typeName, ErrUnrecognizedType)
	}

	obj, err := s.vm.New(constructor)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate object of type %s: %w", typeName, err)
	}

	return newGojaScriptObject(obj, s), nil
}

func (s *GojaScriptSystem) importValue(v interface{}) goja.Value {
	if scriptable, ok := v.(Scriptable); ok {
		so := scriptable.GetScript()
		gso := so.(*GojaScriptObject)
		return gso.obj
	} else {
		return s.vm.ToValue(v)
	}
}

func (s *GojaScriptSystem) exportValue(v goja.Value) interface{} {
	if obj, ok := v.(*goja.Object); ok {
		return obj.Get("native")
	} else {
		return v.Export()
	}
}
