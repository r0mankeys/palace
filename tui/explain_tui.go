package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glam "github.com/charmbracelet/glamour"
)

type model struct {
	viewport viewport.Model
	content  string
	ready    bool
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
			// Initialize viewport when window size is known
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = 0
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	}
	m.viewport, _ = m.viewport.Update(msg)
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}
	return m.viewport.View()
}

func RunExplainTUI() error {
	md, err := os.ReadFile("docs/explain.md")
	if err != nil {
		return fmt.Errorf("could not read explain.md: %w", err)
	}

	// Render Markdown with Glamour to ANSI terminals
	renderer, err := glam.NewTermRenderer(
		glam.WithAutoStyle(),
		glam.WithWordWrap(80),
	)
	if err != nil {
		return fmt.Errorf("could not create renderer: %w", err)
	}

	out, err := renderer.Render(string(md))
	if err != nil {
		return fmt.Errorf("could not render markdown: %w", err)
	}

	m := initialModel()
	m.content = out

	p := tea.NewProgram(m, tea.WithAltScreen()) // <-- THIS is the key to full-screen mode
	return p.Start()
}
