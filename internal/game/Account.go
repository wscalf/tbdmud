package game

import (
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

func (a *Account) AddCharacter(p *PlayerSaveData) {
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

type AccountSaveData struct {
	Login          string
	PasswordHash   []byte
	CharactersData []characterEntryData
}

type characterEntryData struct {
	ID   string
	Name string
}

func (a *Account) GetSaveData() *AccountSaveData {
	data := &AccountSaveData{
		Login:          a.Login,
		PasswordHash:   a.passwordHash,
		CharactersData: make([]characterEntryData, 0, len(a.characters)),
	}

	for _, character := range a.characters {
		data.CharactersData = append(data.CharactersData, characterEntryData{
			ID:   character.id,
			Name: character.name,
		})
	}

	return data
}

func AccountFromSaveData(data *AccountSaveData) *Account {
	account := NewAccount(data.Login)
	account.passwordHash = data.PasswordHash

	charactersData := data.CharactersData
	account.characters = make([]characterEntry, 0, len(charactersData))
	for _, charData := range charactersData {
		account.characters = append(account.characters, characterEntry{
			id:   charData.ID,
			name: charData.Name,
		})
	}

	return account
}
