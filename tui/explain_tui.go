package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	glam "github.com/charmbracelet/glamour"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = gloss.NewStyle().
			Border(gloss.RoundedBorder()).
			BorderForeground(gloss.Color("63")).
			Padding(1, 2)

	titleStyle = gloss.NewStyle().
			BorderStyle(gloss.RoundedBorder()).
			Foreground(gloss.Color("212")).
			Padding(0, 1).
			Bold(true)

	infoStyle = gloss.NewStyle().
			BorderStyle(gloss.RoundedBorder()).
			Foreground(gloss.Color("241")).
			Padding(0, 1)
)

type model struct {
	viewport viewport.Model
	content  string
	ready    bool
	width    int
	height   int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// The bordered “box” will take up ~80% of the screen
		boxWidth := int(float64(msg.Width) * 0.8)
		boxHeight := int(float64(msg.Height) * 0.8)

		// Account for header/footer
		headerHeight := gloss.Height(m.headerView())
		footerHeight := gloss.Height(m.footerView())
		verticalMargin := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(boxWidth-4, boxHeight-verticalMargin-4) // smaller for padding
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = boxWidth - 4
			m.viewport.Height = boxHeight - verticalMargin - 4
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	content := fmt.Sprintf(
		"%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	)

	box := borderStyle.Render(content)

	// Center box horizontally and vertically
	hpad := (m.width - gloss.Width(box)) / 2
	vpad := (m.height - gloss.Height(box)) / 2
	if hpad < 0 {
		hpad = 0
	}
	if vpad < 0 {
		vpad = 0
	}

	return gloss.Place(m.width, m.height, gloss.Center, gloss.Center, box)
}

func (m model) headerView() string {
	title := titleStyle.Render("palace — Explain Mode")
	line := strings.Repeat("─", max(0, m.viewport.Width-gloss.Width(title)))
	return gloss.JoinHorizontal(gloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%  |  ↑↓ scroll  |  q quit", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-gloss.Width(info)))
	return gloss.JoinHorizontal(gloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RunExplainTUI() error {
	md, err := os.ReadFile("docs/explain.md")
	if err != nil {
		return fmt.Errorf("could not read explain.md: %w", err)
	}

	// Render Markdown with Glamour
	renderer, err := glam.NewTermRenderer(
		glam.WithAutoStyle(),
		glam.WithWordWrap(76),
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

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	return p.Start()
}
