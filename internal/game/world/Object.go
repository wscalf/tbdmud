package world

import "github.com/wscalf/tbdmud/internal/game/contracts"

type Object struct {
	ID          string
	Name        string
	Description string
	script      contracts.ScriptObject
}
