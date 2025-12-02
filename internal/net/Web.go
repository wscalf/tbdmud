package net

import (
	"fmt"
	"net/http"

	"github.com/wscalf/tbdmud/internal/game"
)

type Web struct {
	listener *WebSocketListener
}

func NewWeb() *Web {
	return &Web{
		listener: NewWebSocketListener(),
	}
}

func (w *Web) Listener() game.ClientListener {
	return w.listener
}

func (w *Web) ServeHttp(port int) {
	http.HandleFunc("/socket", w.listener.ServeWS)
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}
