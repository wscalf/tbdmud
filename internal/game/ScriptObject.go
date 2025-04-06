package game

type ScriptObject interface {
	Get(prop string) (any, error)
	Set(prop string, value any) error
	Call(name string, args ...any) (any, error)
	GetDescribeProperties() (map[string]any, error)
	GetSaveProperties() (map[string]any, error)
	Type() string
}
