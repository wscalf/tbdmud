package web

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/net"
)

//go:embed index.html
var indexHtml []byte

type Web struct {
	listener *net.WebSocketListener
}

func NewWeb() *Web {
	return &Web{
		listener: net.NewWebSocketListener(),
	}
}

func (w *Web) Listener() game.ClientListener {
	return w.listener
}

func (w *Web) ServeHttp(port int) {
	http.HandleFunc("/socket", w.listener.ServeWS)
	http.HandleFunc("/", w.serveIndex)
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

func (w *Web) serveIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Write(indexHtml)
}
