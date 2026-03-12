package ui

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// RenderMarkdown renders markdown content for the terminal
func RenderMarkdown(content string) string {
	// If content is empty, return empty string
	if strings.TrimSpace(content) == "" {
		return ""
	}

	// Create a new renderer with dark style
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80), // Standard terminal width wrapping
	)
	if err != nil {
		return content // Fallback to raw content if renderer fails
	}

	out, err := r.Render(content)
	if err != nil {
		return content // Fallback to raw content if rendering fails
	}

	return out
}
