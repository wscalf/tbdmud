package scripting

import "github.com/wscalf/tbdmud/internal/game"

type Scriptable interface {
	GetScript() game.ScriptObject
}
