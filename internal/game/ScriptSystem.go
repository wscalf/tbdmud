package game

type ScriptSystem interface {
	RegisterCommands(commands *Commands)
	Wrap(obj interface{}, scriptType string) (ScriptObject, error)
}
