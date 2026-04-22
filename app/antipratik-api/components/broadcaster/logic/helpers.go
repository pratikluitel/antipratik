package logic

import (
	"bytes"

	"github.com/yuin/goldmark"
)

// markdownToHTML converts a markdown string to HTML using goldmark.
func markdownToHTML(md string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return md
	}
	return buf.String()
}
