package game

type Commands struct {
}

func (c *Commands) Execute(p *Player, input string) error {
	return p.Send(input)
}
