package game

type ClientListener interface {
	Listen() (chan Client, error)
	LastError() error
}
