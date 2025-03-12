package parameters

import (
	"strings"
	"unicode"
)

type Delimeter struct {
	text     string
	required bool
}

func (d Delimeter) Name() string {
	return d.text
}

func (d Delimeter) IsRequired() bool {
	return d.required
}

func (d Delimeter) IsMatch(text string) bool {
	start, end := d.find(text)

	if start > -1 && end > -1 {
		return true
	} else {
		return false
	}
}

func (d Delimeter) Consume(text string) (string, string) {
	start, end := d.find(text)

	if start > -1 && end > -1 {
		return "", text[end:]
	} else {
		return "", text
	}
}

func (d Delimeter) find(text string) (int, int) {
	start := 0
	for start < len(text) && unicode.IsSpace(rune(text[start])) {
		start++
	}

	trimmed := text[start:]

	if strings.HasPrefix(trimmed, d.text) {
		return start, start + len(d.text)
	} else {
		return -1, -1
	}
}
