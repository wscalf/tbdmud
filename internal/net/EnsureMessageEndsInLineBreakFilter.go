package net

import "io"

type EnsureMessageEndsInLineBreakFilter struct {
	inner       io.Writer
	lastMessage []byte
}

func NewEnsureMessageEndsInLineBreakFilter(inner io.Writer) *EnsureMessageEndsInLineBreakFilter {
	return &EnsureMessageEndsInLineBreakFilter{
		inner: inner,
	}
}

func (e *EnsureMessageEndsInLineBreakFilter) Write(p []byte) (n int, err error) {
	if e.lastMessage != nil {
		n, err = e.inner.Write(e.lastMessage)
	}

	e.lastMessage = p
	return
}

func (e *EnsureMessageEndsInLineBreakFilter) WriteFinal() (n int, err error) {
	if e.lastMessage != nil {
		l := len(e.lastMessage)
		if l > 2 {
			endBytes := e.lastMessage[l-2:]
			if endBytes[0] == byte('\r') && endBytes[1] == byte('\n') {
				return e.inner.Write(e.lastMessage) //Already ends in a line break
			}
		}

		msg := append(e.lastMessage, []byte("\r\n")...) //Append a line break if missing
		e.lastMessage = nil
		return e.inner.Write(msg)
	}

	return //Nothing to do if there's no last message
}
