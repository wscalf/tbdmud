package net

import (
	"fmt"
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/text"
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

	formatter := text.NewMarkupFilter(writer, handleHTMLFormattingReplacement)
	escaper := NewHTMLEscapeFilter(formatter)

	err = msg.Run(escaper)
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

func handleHTMLFormattingReplacement(directive text.FormattingDirective) string {
	switch directive.FormattingKind {
	case text.FormattingKindForecolor:
		if directive.End {
			return "</span>"
		} else {
			return fmt.Sprintf("<span style=\"color: %s\">", directive.Param) //This may need to be validated to be a color
		}
	case text.FormattingKindBackcolor:
		if directive.End {
			return "</span>"
		} else {
			return fmt.Sprintf("<span style=\"background-color: %s\">", directive.Param) //This may need to be validated to be a color
		}
	case text.FormattingKindBold:
		if directive.End {
			return "</b>"
		} else {
			return "<b>"
		}
	case text.FormattingKindItalic:
		if directive.End {
			return "</i>"
		} else {
			return "<i>"
		}
	case text.FormattingKindUnderline:
		if directive.End {
			return "</u>"
		} else {
			return "<u>"
		}
	case text.FormattingKindStrikeout:
		if directive.End {
			return "</s>"
		} else {
			return "<s>"
		}
	}
	return ""
}
