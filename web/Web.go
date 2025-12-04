package web

import (
	_ "embed"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/net"
)

//go:embed index.html
var indexHtml []byte

type Web struct {
	listener    *net.WebSocketListener
	userContent map[string][]byte
}

func NewWeb(userContent map[string][]byte) *Web {
	return &Web{
		listener:    net.NewWebSocketListener(),
		userContent: userContent,
	}
}

func (w *Web) Listener() game.ClientListener {
	return w.listener
}

func (w *Web) ServeHttp(port int) {
	http.HandleFunc("/socket", w.listener.ServeWS)
	http.HandleFunc("/user-content/", w.serveUserContent)
	http.HandleFunc("/", w.serveIndex)
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
}

func (w *Web) serveUserContent(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	file := path[len(request.Pattern):]

	if content, found := w.userContent[file]; found {
		contentType := mime.TypeByExtension(filepath.Ext(file))
		writer.Header().Add("Content-Type", contentType)
		writer.Write(content)
	} else {
		writer.WriteHeader(404)
	}
}

func (w *Web) serveIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Write(indexHtml)
}
