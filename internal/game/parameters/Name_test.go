package parameters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName_IsMatchAndConsume(t *testing.T) {
	tests := []struct {
		text              string
		expectedToExtract string
		expectedToRemain  string
		expectedToMatch   bool
	}{
		{"", "", "", false},
		{"foo", "foo", "", true},
		{"foo bar", "foo", " bar", true},
		{" foo bar", "foo", " bar", true},
		{`"foo bar"`, "foo bar", "", true},
		{`"foo bar" jar`, "foo bar", " jar", true},
		{` "foo bar" jar`, "foo bar", " jar", true},
		{`"foo`, "foo", "", true},
	}

	name := Name{}

	for _, test := range tests {
		matches := name.IsMatch(test.text)
		assert.Equal(t, test.expectedToMatch, matches, "incorrect match result for: %+v", test)

		extracted, remainder := name.Consume(test.text)
		assert.Equal(t, test.expectedToExtract, extracted, "incorrect value extracted for: %+v", test)
		assert.Equal(t, test.expectedToRemain, remainder, "incorrect remaining text for: %+v", test)
	}
}
