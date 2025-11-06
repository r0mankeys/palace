package tui

import (
	"fmt"
	"os"
	glam "github.com/charmbracelet/glamour"
)

func RenderExplainPage() error {
	content, err := os.ReadFile("docs/explain.md")
	if err != nil {
		return fmt.Errorf("could not read explain.md: %w", err)
	}

	// Use glamour with a dark or light theme
	renderer, err := glam.NewTermRenderer(
		glam.WithAutoStyle(),
		glam.WithWordWrap(80),
	)
	if err != nil {
		return fmt.Errorf("could not create renderer: %w", err)
	}

	out, err := renderer.Render(string(content))
	if err != nil {
		return fmt.Errorf("could not render markdown: %w", err)
	}

	fmt.Println(out)
	return nil
}
