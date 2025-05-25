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

const sgrPattern string = "\033[%sm"

var ansiBold string = fmt.Sprintf(sgrPattern, "1")
var ansiFaint string = fmt.Sprintf(sgrPattern, "2")
var ansiItalic string = fmt.Sprintf(sgrPattern, "3")
var ansiUnderline string = fmt.Sprintf(sgrPattern, "4")
var ansiStrikeout string = fmt.Sprintf(sgrPattern, "9")
var ansiClearBoldFaint string = fmt.Sprintf(sgrPattern, "22")
var ansiClearItalic string = fmt.Sprintf(sgrPattern, "23")
var ansiClearUnderline string = fmt.Sprintf(sgrPattern, "24")
var ansiClearStrikeout string = fmt.Sprintf(sgrPattern, "29")
var ansiForecolorBlack string = fmt.Sprintf(sgrPattern, "30")
var ansiForecolorRed string = fmt.Sprintf(sgrPattern, "31")
var ansiForecolorGreen string = fmt.Sprintf(sgrPattern, "32")
var ansiForecolorYellow string = fmt.Sprintf(sgrPattern, "33")
var ansiForecolorBlue string = fmt.Sprintf(sgrPattern, "34")
var ansiForecolorMagenta string = fmt.Sprintf(sgrPattern, "35")
var ansiForecolorCyan string = fmt.Sprintf(sgrPattern, "36")
var ansiForecolorWhite string = fmt.Sprintf(sgrPattern, "37")
var ansiClearforecolor string = fmt.Sprintf(sgrPattern, "39")
var ansiBackcolorBlack string = fmt.Sprintf(sgrPattern, "40")
var ansiBackcolorRed string = fmt.Sprintf(sgrPattern, "41")
var ansiBackcolorGreen string = fmt.Sprintf(sgrPattern, "42")
var ansiBackcolorYellow string = fmt.Sprintf(sgrPattern, "43")
var ansiBackcolorBlue string = fmt.Sprintf(sgrPattern, "44")
var ansiBackcolorMagenta string = fmt.Sprintf(sgrPattern, "45")
var ansiBackcolorCyan string = fmt.Sprintf(sgrPattern, "46")
var ansiBackcolorWhite string = fmt.Sprintf(sgrPattern, "47")
var ansiClearBackcolor string = fmt.Sprintf(sgrPattern, "49")

func handleANSIFormattingReplacement(directive text.FormattingDirective) string {
	switch directive.FormattingKind {
	case text.FormattingKindForecolor:
		if directive.End {
			return ansiClearforecolor
		}

		switch directive.Param {
		case "black":
			return ansiForecolorBlack
		case "red":
			return ansiForecolorRed
		case "green":
			return ansiForecolorGreen
		case "yellow":
			return ansiForecolorYellow
		case "blue":
			return ansiForecolorBlue
		case "magenta":
			return ansiForecolorMagenta
		case "cyan":
			return ansiForecolorCyan
		case "white":
			return ansiForecolorWhite
		default:
			return ansiForecolorWhite
		}
	case text.FormattingKindBackcolor:
		if directive.End {
			return ansiClearBackcolor
		}

		switch directive.Param {
		case "black":
			return ansiBackcolorBlack
		case "red":
			return ansiBackcolorRed
		case "green":
			return ansiBackcolorGreen
		case "yellow":
			return ansiBackcolorYellow
		case "blue":
			return ansiBackcolorBlue
		case "magenta":
			return ansiBackcolorMagenta
		case "cyan":
			return ansiBackcolorCyan
		case "white":
			return ansiBackcolorWhite
		default:
			return ansiBackcolorBlack
		}

	case text.FormattingKindBold:
		if directive.End {
			return ansiClearBoldFaint
		} else {
			return ansiBold
		}

	case text.FormattingKindFaint:
		if directive.End {
			return ansiClearBoldFaint
		} else {
			return ansiFaint
		}

	case text.FormattingKindItalic:
		if directive.End {
			return ansiClearItalic
		} else {
			return ansiItalic
		}

	case text.FormattingKindUnderline:
		if directive.End {
			return ansiClearUnderline
		} else {
			return ansiUnderline
		}

	case text.FormattingKindStrikeout:
		if directive.End {
			return ansiClearStrikeout
		} else {
			return ansiStrikeout
		}
	}

	return ""
}
