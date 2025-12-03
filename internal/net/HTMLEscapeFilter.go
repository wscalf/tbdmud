package net

import (
	"io"
)

var htmlReplacements = map[rune][]byte{
	'<':  []byte("&lt;"),
	'>':  []byte("&gt;"),
	'&':  []byte("&amp;"),
	'"':  []byte("&#34"),
	'\'': []byte("&#39;"),
}

type HTMLEscapeFilter struct {
	inner io.Writer
}

func NewHTMLEscapeFilter(inner io.Writer) *HTMLEscapeFilter {
	return &HTMLEscapeFilter{
		inner: inner,
	}
}

func (h *HTMLEscapeFilter) Write(p []byte) (int, error) {
	s := string(p)
	written := 0
	start := 0
	for i, c := range s {
		if replacement, matched := htmlReplacements[c]; matched {
			cleanPortion := s[start:i]
			n, err := h.inner.Write([]byte(cleanPortion))
			if err != nil {
				return 0, err
			}
			written += n

			n, err = h.inner.Write(replacement)
			if err != nil {
				return 0, err
			}
			written += n
			start = i + 1 //Skip replaced character
		}
	}

	if start < len(s) {
		cleanPortion := s[start:]
		n, err := h.inner.Write([]byte(cleanPortion))
		if err != nil {
			return 0, err
		}
		written += n
	}

	return written, nil
}
