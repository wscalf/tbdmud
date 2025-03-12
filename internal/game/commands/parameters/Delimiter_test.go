package parameters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelimiter_IsMatch(t *testing.T) {
	cases := []struct {
		text     string
		expected bool
	}{
		{"val", true},
		{"val ", true},
		{" val ", true},
		{" val", true},
		{"mal", false},
		{"mal ", false},
		{" mal ", false},
		{" mal", false},
		{"", false},
	}

	delimiter := &Delimeter{
		text:     "val",
		required: true,
	}

	for _, test := range cases {
		actual := delimiter.IsMatch(test.text)

		assert.Equal(t, test.expected, actual, "%v", test)
	}
}

func TestDelimiter_Consume(t *testing.T) {
	cases := []struct {
		text     string
		leftover string
	}{
		{"val", ""},
		{"val ", " "},
		{" val ", " "},
		{" val", ""},
		{"mal", "mal"},
		{"mal ", "mal "},
		{" mal ", " mal "},
		{" mal", " mal"},
		{"", ""},
	}

	delimiter := &Delimeter{
		text:     "val",
		required: true,
	}

	for _, test := range cases {
		result, leftover := delimiter.Consume(test.text)

		assert.Equal(t, "", result, "Delimiters should not produce a value. Case: %+v", test)
		assert.Equal(t, test.leftover, leftover, "Incorrect leftover text after consuming. Case: %+v", test)
	}
}
