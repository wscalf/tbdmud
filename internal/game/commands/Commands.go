package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wscalf/tbdmud/internal/game/commands/parameters"
	"github.com/wscalf/tbdmud/internal/game/jobs"
	"github.com/wscalf/tbdmud/internal/game/world"
)

var InputError error = errors.New("invalid input")

type Commands struct {
	commands map[string]Command
}

func NewCommands() *Commands {
	return &Commands{
		commands: map[string]Command{},
	}
}

func (c *Commands) RegisterBuiltins() {
	c.Register("think", Think{})
	c.Register("help", Help{commands: c})
}

func (c *Commands) Register(name string, command Command) {
	c.commands[name] = command
}

func (c *Commands) Prepare(p *world.Player, input string) (jobs.Job, error) {
	name, argPart := splitCommandNameFromArgs(input)

	command, ok := c.commands[name]
	if !ok {
		return nil, fmt.Errorf("%w: unrecognized command %s: try help", InputError, name)
	}

	parameterSpec := command.GetParameters()
	parameters, err := extractParameters(name, argPart, parameterSpec)

	if err != nil {
		return nil, err
	}

	return &commandJob{
		command: command,
		player:  p,
		params:  parameters,
	}, nil
}

func extractParameters(cmd string, text string, parameterSpec []parameters.Parameter) (map[string]string, error) {
	args := make(map[string]string)
	var value string

	for _, p := range parameterSpec {
		if p.IsMatch(text) {
			value, text = p.Consume(text)
			args[p.Name()] = value
		} else {
			if p.IsRequired() {
				return nil, fmt.Errorf("%w: missing required parameter %s. Try help %s", InputError, p.Name(), cmd)
			}
		}
	}

	return args, nil
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
