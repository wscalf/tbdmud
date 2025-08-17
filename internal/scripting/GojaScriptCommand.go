package scripting

import (
	"log/slog"

	"github.com/dop251/goja"
	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/game/parameters"
)

type GojaScriptCommand struct {
	name        string
	description string
	parameters  []parameters.Parameter
	system      *GojaScriptSystem
	function    goja.Callable
}

func NewGojaScriptCommand(
	name string,
	description string,
	parameters []parameters.Parameter,
	system *GojaScriptSystem,
	function goja.Callable) *GojaScriptCommand {
	return &GojaScriptCommand{
		name:        name,
		description: description,
		parameters:  parameters,
		system:      system,
		function:    function,
	}
}

func (c *GojaScriptCommand) GetDescription() string {
	return c.description
}

func (c *GojaScriptCommand) GetParameters() []parameters.Parameter {
	return c.parameters
}

func (c *GojaScriptCommand) Execute(player *game.Player, args map[string]string, state map[string]any, requeueHandler func()) bool {
	importedArgs := make([]goja.Value, 0, len(c.parameters)+1) //Need room for all declared parameters + player
	importedArgs = append(importedArgs, c.system.importValue(player))
	for _, parameter := range c.parameters {
		if val, ok := args[parameter.Name()]; ok {
			importedArgs = append(importedArgs, c.system.importValue(val))
		} else {
			importedArgs = append(importedArgs, goja.Null())
		}
	}

	_, err := c.function(goja.Null(), importedArgs...)
	if err != nil {
		slog.Error("error processing soft-coded command", "err", err, "cmd", c.name, "args", args)
	}

	return true
}
