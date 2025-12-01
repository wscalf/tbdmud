package net

import (
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/wscalf/tbdmud/internal/game"
)

type WebSocketClient struct {
	ws         *websocket.Conn
	input      chan string
	lastError  error
	disconnect bool
}

func newWebSocketClient(ws *websocket.Conn) *WebSocketClient {
	return &WebSocketClient{
		ws:    ws,
		input: make(chan string, 5),
	}
}

func (c *WebSocketClient) Send(msg game.OutputJob) error {
	writer, err := c.ws.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	defer writer.Close()

	err = msg.Run(writer)
	return err
}

func (c *WebSocketClient) Recv() chan string {
	return c.input
}

func (c *WebSocketClient) LastError() error {
	return c.lastError
}

func (c *WebSocketClient) Disconnect() {
	c.disconnect = true
	c.ws.Close()
}

func (c *WebSocketClient) Run() {
	defer c.ws.Close()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Debug("websocket client disconnect", "data", err)
			} else {
				slog.Error("websocket receive error", "error", err)
				c.lastError = err
			}
			break
		}
		body := string(message)

		c.input <- body
	}
}
