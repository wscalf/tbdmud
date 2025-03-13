package world

type TransitionJob struct {
	Player *Player
	To     *Room
}

func (t *TransitionJob) Run() {
	if t.Player.room != nil {
		t.Player.Leave()
	}

	t.Player.Join(t.To)
}
