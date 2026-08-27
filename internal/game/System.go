package game

import "time"

type System struct {
}

func NewSystem() *System {
	return &System{}
}

func (s *System) Wait(duration time.Duration) {
	time.Sleep(duration)
}
