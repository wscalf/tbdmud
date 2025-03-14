package game

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/wscalf/tbdmud/internal/game/parameters"
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
	c.Register("say", Say{})
	c.Register("look", Look{})
}

func (c *Commands) Register(name string, command Command) {
	c.commands[name] = command
}

func (c *Commands) Prepare(p *Player, input string) (Job, error) {
	name, argPart := SplitCommandNameFromArgs(input)
	var command Command
	room := p.GetRoom()

	if room == nil {
		slog.Error("player tried to execute a command while not in a room.", "player", p.Name)
		return nil, fmt.Errorf("internal error")
	}
	command = room.FindLocalCommand(name)
	if command == nil {
		command = c.commands[name]
	}

	if command == nil {
		return nil, fmt.Errorf("%w: unrecognized command %s: try help", InputError, name)
	}

	parameterSpec := command.GetParameters()
	parameters, err := ExtractParameters(name, argPart, parameterSpec)

	if err != nil {
		return nil, err
	}

	return &commandJob{
		command: command,
		player:  p,
		params:  parameters,
	}, nil
}

func ExtractParameters(cmd string, text string, parameterSpec []parameters.Parameter) (map[string]string, error) {
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

func SplitCommandNameFromArgs(input string) (string, string) {
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
	player  *Player
	params  map[string]string
}

func (c *commandJob) Run() {
	c.command.Execute(c.player, c.params)
}
