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
	_, found := state["asyncContext"]
	if !found {
		importedArgs := make([]goja.Value, 0, len(c.parameters)+1) //Need room for all declared parameters + player
		importedArgs = append(importedArgs, c.system.importValue(player))
		for _, parameter := range c.parameters {
			if val, ok := args[parameter.Name()]; ok {
				importedArgs = append(importedArgs, c.system.importValue(val))
			} else {
				importedArgs = append(importedArgs, goja.Null())
			}
		}

		asyncContext := NewGojaAsyncContext(requeueHandler)
		c.system.setCurrentCommandAsyncContext(asyncContext)
		defer c.system.clearCurrentCommandAsyncContext()
		state["asyncContext"] = asyncContext

		result, err := c.function(goja.Null(), importedArgs...)
		if err != nil {
			slog.Error("error processing soft-coded command", "err", err, "cmd", c.name, "args", args)
		}

		promise, ok := result.Export().(*goja.Promise)
		if ok {
			state["promise"] = promise
			return false
		} else {
			return true
		}
	} else {
		promise := state["promise"].(*goja.Promise)
		asyncContext := state["asyncContext"].(*GojaAsyncContext)
		err := asyncContext.Resolve()
		if err != nil {
			slog.Error("error on resume", "err", err, "cmd", c.name, "args", args, "context", asyncContext)
			return true
		}
		return isDone(promise)
	}
}
