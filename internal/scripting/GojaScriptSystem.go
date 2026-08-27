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
	vm                         *goja.Runtime
	types                      map[string]goja.Value
	getPersistedProperties     goja.Callable
	currentCommandAsyncContext *GojaAsyncContext
}

func NewGojaScriptSystem() *GojaScriptSystem {
	return &GojaScriptSystem{
		vm:    goja.New(),
		types: map[string]goja.Value{},
	}
}

func (s *GojaScriptSystem) RunBootstrapCode() error {
	err := s.Run(bootstrapCode)
	if err != nil {
		return err
	}

	s.populateKnownTypes()
	return nil
}

func (s *GojaScriptSystem) Run(src string) error {
	_, err := s.vm.RunString(src)

	return err
}

func (s *GojaScriptSystem) Wrap(native interface{}, typeName string) (game.ScriptObject, error) {
	obj, err := s.wrap(native, typeName)
	if err != nil {
		return nil, err
	}

	return newGojaScriptObject(obj, s, typeName), nil
}

func (s *GojaScriptSystem) wrap(native interface{}, typeName string) (*goja.Object, error) {
	native = s.applyDecorator(native)
	obj, err := s.createObject(typeName)
	if err != nil {
		return nil, err
	}
	obj.Set("native", native) //Bypass the setter for native pointer so it doesn't try to unwrap it
	return obj, nil
}

func (s *GojaScriptSystem) applyDecorator(native any) any {
	switch n := native.(type) {
	case *game.Players:
		return &PlayersWrapper{players: n, system: s}
	case *game.System:
		return &SystemWrapper{system: s, sys: n}
	default:
		return native
	}
}

func (s *GojaScriptSystem) createPromise() *GojaPromiseWrapper {
	promise, resolve, reject := s.vm.NewPromise()
	return NewGojaPromiseWrapper(promise, resolve, reject)
}

func (s *GojaScriptSystem) AddGlobal(name, scriptType string, native interface{}) error {
	scriptObj, err := s.wrap(native, scriptType)
	if err != nil {
		return err
	}

	return s.vm.Set(name, scriptObj)
}

func (s *GojaScriptSystem) populateKnownTypes() {
	for _, name := range s.vm.GlobalObject().GetOwnPropertyNames() {
		if _, found := s.types[name]; found { //Skip already known types
			continue
		}

		candidate := s.vm.Get(name)
		if _, ok := goja.AssertConstructor(candidate); ok {
			s.types[name] = candidate
		}
	}
}

func (s *GojaScriptSystem) Initialize() error {
	//Find all the types
	s.populateKnownTypes()

	f, ok := goja.AssertFunction(s.vm.Get("getPersistedProperties"))
	if !ok {
		return fmt.Errorf("getPersistedProperties not found in bootstrap code")
	}

	s.getPersistedProperties = f
	return nil
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

func (s *GojaScriptSystem) getPersistedFieldsForType(typeName string) ([]string, error) {
	v, err := s.getPersistedProperties(goja.Null(), s.importValue(typeName))
	if err != nil {
		return nil, err
	}

	values := v.Export().([]any)
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.(string))
	}

	return result, nil
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

func (s *GojaScriptSystem) createObject(typeName string) (*goja.Object, error) {
	constructor, ok := s.types[typeName]
	if !ok {
		return nil, fmt.Errorf("unable to instantiate object of type %s: %w", typeName, ErrUnrecognizedType)
	}

	obj, err := s.vm.New(constructor)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate object of type %s: %w", typeName, err)
	}

	return obj, nil
}

func (s *GojaScriptSystem) importValue(v interface{}) goja.Value {
	if scriptable, ok := v.(Scriptable); ok {
		so := scriptable.GetScript()
		gso := so.(*GojaScriptObject)
		return gso.Obj
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

func (s *GojaScriptSystem) setCurrentCommandAsyncContext(ctx *GojaAsyncContext) {
	s.currentCommandAsyncContext = ctx
}

func (s *GojaScriptSystem) getCurrentCommandAsyncContext() *GojaAsyncContext {
	return s.currentCommandAsyncContext
}

func (s *GojaScriptSystem) clearCurrentCommandAsyncContext() {
	s.currentCommandAsyncContext = nil
}
