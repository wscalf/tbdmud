package game

import (
	"github.com/wscalf/tbdmud/internal/game/parameters"
	"github.com/wscalf/tbdmud/internal/text"
)

/*
Think
*/
var thinkparams = []parameters.Parameter{parameters.NewFreeText("thought")}

type Think struct{}

func (t Think) GetDescription() string {
	return "Sends the thought text to the player."
}

func (t Think) GetParameters() []parameters.Parameter {
	return thinkparams
}

func (t Think) Execute(player *Player, args map[string]string) {
	thought := args["thought"]

	player.Sendf(thought)
}

/*
Say
*/
var sayparams []parameters.Parameter = []parameters.Parameter{parameters.NewFreeText("text")}

type Say struct{}

func (s Say) GetDescription() string {
	return "Sends the thought text to the player."
}

func (s Say) GetParameters() []parameters.Parameter {
	return sayparams
}

func (s Say) Execute(player *Player, args map[string]string) {
	text := args["text"]
	player.Sendf(`You say, "%s"`, text)

	room := player.GetRoom()
	room.SendToAllExcept(player, `%s says, "%s"`, player.Name, text)
}

/*
Look
*/
var lookparams []parameters.Parameter = []parameters.Parameter{parameters.NewName("object", false)}

type Look struct{}

func (l Look) GetDescription() string {
	return "Describes the current room or the object looked at."
}

func (l Look) GetParameters() []parameters.Parameter {
	return lookparams
}

func (l Look) Execute(player *Player, args map[string]string) {
	name, found := args["object"]
	room := player.GetRoom()

	if !found {
		player.Send(room.Describe())
		return
	}

	other := room.FindPlayer(name)
	if other != nil {
		player.Send(other.Describe())
		other.Sendf("%s looked at you", player.Name)
		return
	}

	item := room.FindItem(name)
	if item != nil {
		player.Sendf(item.Description) //TODO make this use a layout
		return
	}

	player.Sendf("I don't see that here.")
}

/*
Help
*/
var helpparams []parameters.Parameter = []parameters.Parameter{parameters.NewName("cmd", false)}

type Help struct {
	commands *Commands
}

func (h Help) GetDescription() string {
	return "Prints this help text."
}

func (h Help) GetParameters() []parameters.Parameter {
	return helpparams
}

func (h Help) Execute(player *Player, params map[string]string) {
	if name, ok := params["cmd"]; ok {
		//A specific command was passed
		if cmd, ok := h.commands.commands[name]; ok {
			params := cmd.GetParameters()
			player.Sendf("%s: %s", name, cmd.GetDescription())
			usage := "Usage: " + name
			for _, param := range params {
				usage = usage + " [" + param.Name() + "]"
			}
			player.Sendf(usage)
		} else {
			player.Sendf("Unrecognized command: %s", name)
		}
	} else {
		//No command was passed, print the list
		//This should probably be externalized to a template
		player.Sendf("The following commands are available:")
		for name, cmd := range h.commands.commands {
			player.Sendf("%s: %s", name, cmd.GetDescription())
		}
	}
}

/*
Take
*/
var takeparams = []parameters.Parameter{parameters.NewName("item", true)}

type Take struct{}

func (t Take) GetDescription() string {
	return "Picks up <item> from the room."
}

func (t Take) GetParameters() []parameters.Parameter {
	return takeparams
}

func (t Take) Execute(player *Player, args map[string]string) {
	itemName := args["item"]
	room := player.GetRoom()

	item := room.FindItem(itemName)
	if item == nil {
		player.Sendf("There's no %s here.", itemName)
		return
	}

	room.RemoveItem(item)
	player.Give(item)

	player.Sendf("You pick up the %s.", itemName)
	room.SendToAllExcept(player, "%s picks up the %s.", player.Name, itemName)
}

/*
Give
*/
var giveparams = []parameters.Parameter{parameters.NewName("item", true), parameters.NewName("player", true)}

type Give struct{}

func (t Give) GetDescription() string {
	return "Gives <item> to <player>."
}

func (t Give) GetParameters() []parameters.Parameter {
	return giveparams
}

func (t Give) Execute(player *Player, args map[string]string) {
	itemName := args["item"]
	toName := args["player"]

	item := player.FindItem(itemName)
	if item == nil {
		player.Sendf("You don't have a %s.", itemName)
		return
	}

	to := player.GetRoom().FindPlayer(toName)
	if to == nil {
		player.Sendf("There is no %s here.", toName)
		return
	}

	player.Take(item)
	to.Give(item)

	player.Sendf("You give the %s to %s.", itemName, toName)
	to.Sendf("%s gives you the %s.", player.Name, itemName)
}

/*
Inv
*/
var invparams = []parameters.Parameter{}

type Inv struct {
	layout *text.Layout
}

func (i Inv) GetDescription() string {
	return "Lists the items you have."
}

func (i Inv) GetParameters() []parameters.Parameter {
	return invparams
}

func (i Inv) Execute(player *Player, args map[string]string) {
	job := i.layout.Prepare(player)
	player.Send(job)
}

/*
Inv
*/
var descparams = []parameters.Parameter{parameters.NewFreeText("description")}

type Desc struct{}

func (d Desc) GetDescription() string {
	return "Changes your character's description."
}

func (d Desc) GetParameters() []parameters.Parameter {
	return descparams
}

func (d Desc) Execute(player *Player, args map[string]string) {
	desc := args["description"]
	player.Description = desc
}

func (c *Commands) RegisterBuiltins(layouts map[string]*text.Layout) {
	c.Register("think", Think{})
	c.Register("desc", Desc{})
	c.Register("help", Help{commands: c})
	c.Register("say", Say{})
	c.Register("look", Look{})
	c.Register("take", Take{})
	c.Register("give", Give{})
	c.Register("inv", Inv{layout: layouts["inventory"]})
}
