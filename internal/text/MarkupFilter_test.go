package text

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkupFilterPassesThroughUnmarkedText(t *testing.T) {
	sb := &strings.Builder{}
	filter := NewMarkupFilter(sb, nil)

	input := "Foo"
	filter.Write([]byte(input))
	output := sb.String()

	assert.Equal(t, input, output)
}
func TestMarkupFilterSkipsUnhandledDirectives(t *testing.T) {
	sb := &strings.Builder{}
	filter := NewMarkupFilter(sb, nil)

	input := "[fc=red]F[i]o[/i]o[/fc]"
	filter.Write([]byte(input))
	output := sb.String()

	assert.Equal(t, "Foo", output)
}

func TestMarkupFilterHandlesReplacements(t *testing.T) {
	sb := &strings.Builder{}
	filter := NewMarkupFilter(sb, handleReplacement)

	input := "[fc=red]F[i]o[/i]o[/fc]"
	filter.Write([]byte(input))
	output := sb.String()

	assert.Equal(t, "(red)F(i)o(!i)o(!color)", output)
}

func TestMarkupFilterHandlesDirectivesSplitOverWriteCalls(t *testing.T) {
	sb := &strings.Builder{}
	filter := NewMarkupFilter(sb, handleReplacement)

	// Original input: "[fc=red]F[i]o[/i]o[/fc]"
	filter.Write([]byte("["))
	filter.Write([]byte("fc=red"))
	filter.Write([]byte("]F[i"))
	filter.Write([]byte("]o[/i]o["))
	filter.Write([]byte("/"))
	filter.Write([]byte("f"))
	filter.Write([]byte("c"))
	filter.Write([]byte("]"))
	output := sb.String()

	assert.Equal(t, "(red)F(i)o(!i)o(!color)", output)
}

func handleReplacement(directive FormattingDirective) string {
	switch directive.FormattingKind {
	case FormattingKindForecolor:
		if directive.End {
			return "(!color)"
		} else {
			return fmt.Sprintf("(%s)", directive.Param)
		}
	case FormattingKindItalic:
		if directive.End {
			return "(!i)"
		} else {
			return "(i)"
		}
	default:
		return ""
	}
}
