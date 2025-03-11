package contracts

type ClientListener interface {
	Listen() (chan Client, error)
	LastError() error
}
