package net

import (
	"errors"
	"io"
	"strings"

	"github.com/reiver/go-telnet"
)

const (
	CR byte = 13
	LF byte = 10
)

type TelnetClient struct {
	reader     telnet.Reader
	writer     telnet.Writer
	input      chan string
	disconnect bool
	lastError  error
}

func newTelnetClient(reader telnet.Reader, writer telnet.Writer) *TelnetClient {
	return &TelnetClient{
		reader: reader,
		writer: writer,
	}
}

func (t *TelnetClient) Send(msg string) error {
	_, err := t.writer.Write([]byte(msg))
	return err
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
