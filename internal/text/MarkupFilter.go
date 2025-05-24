package text

import (
	"io"
	"strings"
)

const (
	FormattingKindUnknown = iota
	FormattingKindBold
	FormattingKindFaint
	FormattingKindItalic
	FormattingKindUnderline
	FormattingKindStrikeout
	FormattingKindForecolor
	FormattingKindBackcolor
)

type FormattingDirective struct {
	FormattingKind int
	End            bool
	Param          string
}

func parseFormattingDirective(raw string) FormattingDirective {
	start := 0
	directive := FormattingDirective{}
	if raw[start] == '/' {
		directive.End = true
		start += 1
	}

	equal := strings.IndexByte(raw[start:], '=')
	if equal == -1 {
		directive.FormattingKind = parseFormattingKind(raw[start:])
	} else {
		directive.FormattingKind = parseFormattingKind(raw[start:equal])
		directive.Param = raw[equal+1:]
	}

	return directive
}

func parseFormattingKind(raw string) int {
	switch raw {
	case "b":
		return FormattingKindBold
	case "f":
		return FormattingKindFaint
	case "i":
		return FormattingKindItalic
	case "u":
		return FormattingKindUnderline
	case "s":
		return FormattingKindStrikeout
	case "fc":
		return FormattingKindForecolor
	case "bc":
		return FormattingKindBackcolor
	default:
		return FormattingKindUnknown
	}
}

type MarkupFilter struct {
	partial             *strings.Builder
	directiveInProgress bool
	inner               io.Writer
	directiveCallback   func(FormattingDirective) string
}

func NewMarkupFilter(inner io.Writer, directiveCallback func(FormattingDirective) string) *MarkupFilter {
	return &MarkupFilter{
		partial:             &strings.Builder{},
		directiveInProgress: false,
		inner:               inner,
		directiveCallback:   directiveCallback,
	}
}

func (m *MarkupFilter) Write(p []byte) (n int, err error) {
	str := string(p)

	start := 0

	if m.directiveInProgress {
		close := strings.IndexByte(str, ']')
		if close == -1 { //Most likely the directive is spanning write calls - probably need to add to an internal buffer
			m.partial.WriteString(str)
			return 0, nil
		}

		m.partial.WriteString(str[:close])
		start = close + 1

		if m.directiveCallback != nil {
			raw := m.partial.String()
			m.partial.Reset()
			directive := parseFormattingDirective(raw)
			replacement := m.directiveCallback(directive)

			written, err := m.inner.Write([]byte(replacement))
			n += written
			if err != nil {
				return n, err
			}
		}
	}

	for {
		open := strings.IndexByte(str[start:], '[')
		if open == -1 {
			written, err := m.inner.Write([]byte(str[start:])) //If there are no more open directives, send the rest of the text on. May need a check that start < len
			n += written
			if err != nil {
				return n, err
			}

			break
		}

		open += start

		if open > start { //There's text before the next open directive
			written, err := m.inner.Write([]byte(str[start:open]))
			n += written
			if err != nil {
				return n, err
			}
		}

		close := strings.IndexByte(str[open:], ']')
		if close == -1 { //Most likely the directive is spanning write calls - probably need to add to an internal buffer
			m.partial.WriteString(str[open+1:])
			m.directiveInProgress = true
			break
		}

		close += open
		if m.directiveCallback != nil {
			raw := str[open+1 : close]
			directive := parseFormattingDirective(raw)
			replacement := m.directiveCallback(directive)

			written, err := m.inner.Write([]byte(replacement))
			n += written
			if err != nil {
				return n, err
			}
		}
		start = close + 1
	}

	return n, nil
}
