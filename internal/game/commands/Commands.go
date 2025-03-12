package commands

import (
	"strings"

	"github.com/wscalf/tbdmud/internal/game/commands/builtins"
	"github.com/wscalf/tbdmud/internal/game/commands/parameters"
	"github.com/wscalf/tbdmud/internal/game/jobs"
	"github.com/wscalf/tbdmud/internal/game/world"
)

type Commands struct {
	commands map[string]Command
}

func NewCommands() *Commands {
	return &Commands{
		commands: map[string]Command{},
	}
}

func (c *Commands) RegisterBuiltins() {
	c.Register("think", builtins.Think{})
}

func (c *Commands) Register(name string, command Command) {
	c.commands[name] = command
}

func (c *Commands) Prepare(p *world.Player, input string) jobs.Job {
	name, argPart := splitCommandNameFromArgs(input)

	command, ok := c.commands[name]
	if !ok {
		//Handle unrecognized command
		return nil
	}

	parameterSpec := command.GetParameters()
	parameters := extractParameters(argPart, parameterSpec)

	return &commandJob{
		command: command,
		player:  p,
		params:  parameters,
	}
}

func extractParameters(text string, parameterSpec []parameters.Parameter) map[string]string {
	args := make(map[string]string)
	var value string

	for _, p := range parameterSpec {
		if p.IsMatch(text) {
			value, text = p.Consume(text)
			args[p.Name()] = value
		} else {
			if p.IsRequired() {
				//Handle missing required parameter
			}
		}
	}

	return args
}

func splitCommandNameFromArgs(input string) (string, string) {
	firstSpace := strings.Index(input, " ")
	if firstSpace < 0 {
		return input, "" //It's just the command name
	}

	name := input[0:firstSpace]
	args := input[firstSpace:]

	return strings.TrimSpace(name), strings.TrimSpace(args)
}

type commandJob struct {
	command Command
	player  *world.Player
	params  map[string]string
}

func (c *commandJob) Run() {
	c.command.Execute(c.player, c.params)
}
