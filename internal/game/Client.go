package game

type Client interface {
	Send(msg OutputJob) error
	Recv() chan string
	LastError() error
	Disconnect()
}
