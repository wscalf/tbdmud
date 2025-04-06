package storage

import "github.com/wscalf/tbdmud/internal/game"

type MemoryStore struct {
	accounts   map[string]*game.AccountSaveData
	characters map[string]*game.PlayerSaveData
}

func (s *MemoryStore) CreateOrUpdateAccount(account *game.Account) error {
	s.accounts[account.Login] = account.GetSaveData()
	return nil
}

func (s *MemoryStore) FindAccount(name string) (*game.Account, error) {
	if data, ok := s.accounts[name]; ok {
		return game.AccountFromSaveData(data), nil
	} else {
		return nil, nil
	}
}

func (s *MemoryStore) CreateOrUpdatePlayer(player *game.PlayerSaveData) error {
	s.characters[player.ID] = player
	return nil
}

func (s *MemoryStore) FindPlayer(id string) (*game.PlayerSaveData, error) {
	return s.characters[id], nil
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		accounts:   map[string]*game.AccountSaveData{},
		characters: map[string]*game.PlayerSaveData{},
	}
}
