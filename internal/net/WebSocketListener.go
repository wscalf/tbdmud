package net

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wscalf/tbdmud/internal/game"
)

var upgrader = websocket.Upgrader{}

type WebSocketListener struct {
	clients chan game.Client
}

func NewWebSocketListener(port int) *WebSocketListener {
	return &WebSocketListener{
		clients: make(chan game.Client, 5),
	}
}

func (l *WebSocketListener) Listen() (chan game.Client, error) {
	return l.clients, nil
}

func (l *WebSocketListener) LastError() error {
	return nil
}

func (l *WebSocketListener) ServeWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Error("error upgrading websocket connection", "error", err)
		return
	}

	client := newWebSocketClient(ws)
	l.clients <- client
}
