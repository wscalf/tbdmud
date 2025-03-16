package game

type ScriptObject interface {
	Get(prop string) (interface{}, error)
	Set(prop string, value interface{}) error
	Call(name string, args ...interface{}) (interface{}, error)
}
