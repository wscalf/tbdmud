package game

type Storage interface {
	CreateOrUpdateAccount(account *Account) error
	FindAccount(name string) (*Account, error)
	CreateOrUpdatePlayer(data *PlayerSaveData) error
	FindPlayer(id string) (*PlayerSaveData, error)
}
