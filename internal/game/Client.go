package game

type Client interface {
	Send(msg string) error
	Recv() chan string
	LastError() error
	Disconnect()
}
