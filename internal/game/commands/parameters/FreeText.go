package parameters

import (
	"strings"
)

type FreeText struct {
	name string
}

func NewFreeText(name string) FreeText {
	return FreeText{name: name}
}

func (f FreeText) Name() string {
	return f.name
}

func (f FreeText) IsRequired() bool {
	return false
}

func (f FreeText) IsMatch(text string) bool {
	return true
}

func (f FreeText) Consume(text string) (string, string) {
	result := strings.Trim(text, " ")

	return result, ""
}
