package parameters

import "strings"

type Name struct {
	name     string
	required bool
}

func NewName(name string, required bool) Name {
	return Name{name: name, required: required}
}

func (n Name) Name() string {
	return n.name
}

func (n Name) IsRequired() bool {
	return n.required
}

func (n Name) IsMatch(text string) bool {
	start, end, _ := find(text)
	if start > -1 && end > -1 {
		return true
	} else {
		return false
	}
}

func (n Name) Consume(text string) (string, string) {
	start, end, final := find(text)

	if start < 0 || end < 0 {
		return "", text
	}

	value := text[start:end]
	if len(text) > final {
		return value, text[final:]
	} else {
		return value, ""
	}
}

func find(text string) (int, int, int) {
	delimiter := ""
	start := 0
	for start < len(text) {
		ch := text[start]
		if ch == byte(' ') {
			start++
			continue
		} else if ch == byte('"') {
			//It's a quoted name
			start++ //Skip the quote
			delimiter = "\""
			break
		} else {
			delimiter = " "
			break
		}
	}

	if delimiter == "" {
		return -1, -1, -1
	}

	trimmed := text[start:]
	end := strings.Index(trimmed, delimiter)
	if end < 0 {
		end = len(trimmed)
	}

	if delimiter == "\"" {
		return start, start + end, start + end + 1
	} else {
		return start, start + end, start + end
	}
}
