package game

type ScriptSystem interface {
	RegisterCommands(commands *Commands)
	Wrap(obj interface{}, scriptType string) (ScriptObject, error)
	Run(script string) error
	Initialize() error
	AddGlobal(name, scriptType string, native interface{}) error
}
