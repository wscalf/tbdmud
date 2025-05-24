package net

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/reiver/go-telnet"
	"github.com/wscalf/tbdmud/internal/game"
	"github.com/wscalf/tbdmud/internal/text"
)

const (
	CR  byte = 13
	LF  byte = 10
	ESC byte = 27
)

type TelnetClient struct {
	reader     telnet.Reader
	writer     *text.MarkupFilter
	input      chan string
	disconnect bool
	lastError  error
}

func newTelnetClient(reader telnet.Reader, writer telnet.Writer) *TelnetClient {
	return &TelnetClient{
		reader: reader,
		writer: text.NewMarkupFilter(writer, handleANSIFormattingReplacement),
	}
}

func (t *TelnetClient) Send(msg game.OutputJob) error {
	return msg.Run(t.writer)
}

func (t *TelnetClient) Recv() chan string {
	return t.input
}

func (t *TelnetClient) LastError() error {
	return t.lastError
}

func (t *TelnetClient) Run() {
	var sb strings.Builder
	var buffer [1]byte //Need to read one byte at a time
	t.input = make(chan string)

	for {
		if t.disconnect {
			break //TODO: this currently won't disconnect until a byte is read, but if we expose the underlying TCP connection, we can force it to close eagerly
		}

		n, err := t.reader.Read(buffer[:])
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			if !errors.Is(err, io.EOF) {
				t.lastError = err
			}
			break
		}

		ch := buffer[0]

		switch ch {
		case CR:
			//This is the beginning of a CRLF cluster and the end of a command
			cmd := sb.String()
			sb.Reset()
			t.input <- cmd
		case LF:
			//Ignore LF- they always follow CR
		default:
			sb.WriteByte(buffer[0]) //Append all other data characters to the running buffer
		}
	}

	close(t.input)
}

func (t *TelnetClient) Disconnect() {
	t.disconnect = true
}

const (
	ansiBold             string = "1"
	ansiFaint            string = "2"
	ansiItalic           string = "3"
	ansiUnderline        string = "4"
	ansiStrikeout        string = "9"
	ansiClearBoldFaint   string = "22"
	ansiClearItalic      string = "23"
	ansiClearUnderline   string = "24"
	ansiClearStrikeout   string = "29"
	ansiForecolorBlack   string = "30"
	ansiForecolorRed     string = "31"
	ansiForecolorGreen   string = "32"
	ansiForecolorYellow  string = "33"
	ansiForecolorBlue    string = "34"
	ansiForecolorMagenta string = "35"
	ansiForecolorCyan    string = "36"
	ansiForecolorWhite   string = "37"
	ansiClearforecolor   string = "39"
	ansiBackcolorBlack   string = "40"
	ansiBackcolorRed     string = "41"
	ansiBackcolorGreen   string = "42"
	ansiBackcolorYellow  string = "43"
	ansiBackcolorBlue    string = "44"
	ansiBackcolorMagenta string = "45"
	ansiBackcolorCyan    string = "46"
	ansiBackcolorWhite   string = "47"
	ansiClearBackcolor   string = "49"
)

const sgrPattern string = "\033[%sm"

func handleANSIFormattingReplacement(directive text.FormattingDirective) string {
	switch directive.FormattingKind {
	case text.FormattingKindForecolor:
		if directive.End {
			return fmt.Sprintf(sgrPattern, ansiClearforecolor)
		}

		var code string
		switch directive.Param {
		case "black":
			code = ansiForecolorBlack
		case "red":
			code = ansiForecolorRed
		case "green":
			code = ansiForecolorGreen
		case "yellow":
			code = ansiForecolorYellow
		case "blue":
			code = ansiForecolorBlue
		case "magenta":
			code = ansiForecolorMagenta
		case "cyan":
			code = ansiForecolorCyan
		case "white":
			code = ansiForecolorWhite
		default:
			code = ansiForecolorWhite
		}

		return fmt.Sprintf(sgrPattern, code)
	case text.FormattingKindBackcolor:
		if directive.End {
			return fmt.Sprintf(sgrPattern, ansiClearBackcolor)
		}

		var code string
		switch directive.Param {
		case "black":
			code = ansiBackcolorBlack
		case "red":
			code = ansiBackcolorRed
		case "green":
			code = ansiBackcolorGreen
		case "yellow":
			code = ansiBackcolorYellow
		case "blue":
			code = ansiBackcolorBlue
		case "magenta":
			code = ansiBackcolorMagenta
		case "cyan":
			code = ansiBackcolorCyan
		case "white":
			code = ansiBackcolorWhite
		default:
			code = ansiBackcolorBlack
		}

		return fmt.Sprintf(sgrPattern, code)

	case text.FormattingKindBold:
		if directive.End {
			return fmt.Sprintf(sgrPattern, ansiClearBoldFaint)
		} else {
			return fmt.Sprintf(sgrPattern, ansiBold)
		}

	case text.FormattingKindFaint:
		if directive.End {
			return fmt.Sprintf(sgrPattern, ansiClearBoldFaint)
		} else {
			return fmt.Sprintf(sgrPattern, ansiFaint)
		}

	case text.FormattingKindUnderline:
		if directive.End {
			return fmt.Sprintf(sgrPattern, ansiClearUnderline)
		} else {
			return fmt.Sprintf(sgrPattern, ansiUnderline)
		}

	case text.FormattingKindStrikeout:
		if directive.End {
			return fmt.Sprintf(sgrPattern, ansiClearStrikeout)
		} else {
			return fmt.Sprintf(sgrPattern, ansiStrikeout)
		}
	}

	return ""
}
