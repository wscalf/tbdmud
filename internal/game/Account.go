package game

import (
	"github.com/wscalf/tbdmud/internal/game/world"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Login        string
	passwordHash []byte
	characters   []characterEntry
}

func NewAccount(login string) *Account {
	return &Account{
		Login:      login,
		characters: []characterEntry{},
	}
}

func (a *Account) AddCharacter(p *world.Player) {
	entry := characterEntry{
		id:   p.ID,
		name: p.Name,
	}

	a.characters = append(a.characters, entry)
}

func (a *Account) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.passwordHash = hash
	return nil
}

func (a *Account) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(a.passwordHash, []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

type characterEntry struct {
	id   string
	name string
}
