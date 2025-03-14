package net

import (
	"fmt"

	"github.com/reiver/go-telnet"
	"github.com/wscalf/tbdmud/internal/game"
)

type TelnetListener struct {
	telnet.Server
	clients   chan game.Client
	lastError error
}

func NewTelnetListener(port int) *TelnetListener {
	listener := &TelnetListener{
		Server: telnet.Server{
			Addr:      fmt.Sprintf("localhost:%d", port),
			TLSConfig: nil,
			Logger:    nil,
		},
		clients: make(chan game.Client),
	}

	listener.Handler = listener

	return listener
}

func (t *TelnetListener) Listen() (chan game.Client, error) {
	t.clients = make(chan game.Client)
	go t.ListenAndServe()
	return t.clients, nil
}

func (t *TelnetListener) LastError() error {
	return t.lastError
}

func (t *TelnetListener) ServeTELNET(ctx telnet.Context, writer telnet.Writer, reader telnet.Reader) {
	client := newTelnetClient(reader, writer)
	t.clients <- client
	client.Run()
}
