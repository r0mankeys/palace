// Package tui provides terminal UIs for the palace CLI using the Bubble Tea
// framework
// File `explain_tui.go` displays a centered, styled, scrollable markdown file
// that explains the palace CLI
package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea" 

	// Below are extra styling libraries (purely aesthetic)

	// Viewport: one of the Bubble Tea "bubbles" components for creating a custom
	// window to house the explain markdown file
	"github.com/charmbracelet/bubbles/viewport"
	// Glamour: a markdown renderer for CLI apps 
	glam "github.com/charmbracelet/glamour"
	// Lipgloss: styling utils for terminal layouts
	gloss "github.com/charmbracelet/lipgloss"
)

const (
	// ANSI colours
	MAIN_BORDER_COLOUR = "63" // Cornflower blue
	TITLE_COLOUR = "212" // Lavendar rose
	META_INFO_COLOUR = "241" // Granite grey
	// Meta information about the percentage of the file scrolled, how to escape 
	// etc.
	
	// Magic numbers
	// The documentation isn't clear on what the unit of measure is so I assume
	// it's similar to the CSS `ch` unit, the documentation says "cells"
	PADDING_SMALL = 0
	PADDING_MEDIUM = 1
	PADDING_LARGE = 2
	VIEWPORT_SCALE_FACTOR = 0.8 // ~80% of the screen (subject to change)
	VIEWPORT_PADDING = 4  

	// I/O 
	FILE_NAME = "explain.md"
	FILE_PATH = "docs/" + FILE_NAME
	WORD_WRAP = 76
)

var (
	borderStyle = func() gloss.Style {
		b := gloss.RoundedBorder()
		colour := gloss.Color(MAIN_BORDER_COLOUR)
		return gloss.NewStyle().BorderStyle(b).BorderForeground(colour).Padding(PADDING_MEDIUM, PADDING_LARGE)
	}()

	titleStyle = func() gloss.Style {
		b := gloss.RoundedBorder()
		colour := gloss.Color(TITLE_COLOUR)
		return gloss.NewStyle().BorderStyle(b).Foreground(colour).Padding(PADDING_SMALL, PADDING_MEDIUM).Bold(true)
	}()

	metaInfoStyle = func() gloss.Style {
		b := gloss.RoundedBorder()
		colour := gloss.Color(META_INFO_COLOUR)
		return gloss.NewStyle().BorderStyle(b).Foreground(colour).Padding(PADDING_SMALL, PADDING_MEDIUM)
	}()
)

// The model contains the programs state as well as it's core function, it must
// implement the Model interface so Bubble Tea's internals can create the TUI
// Futher detail: https://pkg.go.dev/github.com/charmbracelet/bubbletea/v2#Model
type model struct {
	viewport viewport.Model
	content  string
	ready    bool // true if the user has already setup a viewport i.e. run the command
	// and may want to resize 
	width    int
	height   int
}

func (m model) Init() tea.Cmd {
	return nil
	// The interface must be fullfiled however we don't need an init command
}

// This function is called whenever a message is received from the user,
// messages contain data from the result of an I/O operation, hence changing the
// UI, essentially think of a message as an event
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// The only messages I care to update as of now are the user wanting to exit
	// the view and the user resizing their terminal window
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	// A Cmd is an I/O function that returns a message

	switch msg := msg.(type) {
	// Handle user wanting to exit 
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Handle user resizing their terminal window
		m.width = msg.Width
		m.height = msg.Height

		vw, vh := calculateViewportDimensions(m.width, m.height)

		// Add margins for the header and footer 
		verticalMargin :=  calculateVerticalMargin(m)

		if !m.ready {
			// `!m.ready` means the user just ran the command
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(vw - VIEWPORT_PADDING, vh - verticalMargin - VIEWPORT_PADDING) // smaller for padding
			m.viewport.YPosition = gloss.Height(m.headerView())
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			// No need to init a new viewport, content etc. just resize
			m.viewport.Width = vw - VIEWPORT_PADDING
			m.viewport.Height = vh - verticalMargin - VIEWPORT_PADDING
		}
	}

	// Finally update the viewport with the message
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View represents a terminal view that can be composed of multiple 
// layers. It can also contain a cursor that will be rendered 
// on top of the layers
func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Build content
	content := fmt.Sprintf(
		"%s\n%s\n%s",
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	)

	contentBox := borderStyle.Render(content)

	// Center box horizontally and vertically
	return centerBox(m, contentBox)
}

// Hepler functions

func calculateViewportDimensions(width, height int) (int, int) {
	vw := int(float64(width) * VIEWPORT_SCALE_FACTOR)
	vh := int(float64(height) * VIEWPORT_SCALE_FACTOR)
	return vw, vh
}

func calculateVerticalMargin(m model) int {
	// Add the heights of the respecitve computed views for a given model
	return gloss.Height(m.headerView()) + gloss.Height(m.footerView())
}

func centerBox(m model, box string) string {
	return gloss.Place(m.width, m.height, gloss.Center, gloss.Center, box)
}

// Styles for the header (showing the title)
func (m model) headerView() string {
	title := titleStyle.Render("palace — Explain Mode")
	line := strings.Repeat("─", getMaxValue(0, m.viewport.Width-gloss.Width(title)))
	return gloss.JoinHorizontal(gloss.Center, title, line)
}

// Styles for the footer (showing the meta data)
func (m model) footerView() string {
	info := metaInfoStyle.Render(fmt.Sprintf("%3.f%%  |  ↑↓ scroll  |  q quit", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", getMaxValue(0, m.viewport.Width-gloss.Width(info)))
	return gloss.JoinHorizontal(gloss.Center, line, info)
}

func getMaxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func glamouriseMarkdown(file []byte) (output string, err error) {
		renderer, err := glam.NewTermRenderer(
		glam.WithAutoStyle(),
		glam.WithWordWrap(WORD_WRAP),
	)
	if err != nil {
		return "", fmt.Errorf("could not create renderer: %w", err)
	}

	out, err := renderer.Render(string(file))
	if err != nil {
		return "", fmt.Errorf("could not render markdown: %w", err)
	}

	return out, nil
}

func loadExplainMarkdown() (file []byte, err error) {
	md, err := os.ReadFile(FILE_PATH)
	if err != nil {
		errMsg := fmt.Sprintf("could not read %s", FILE_NAME)
		return make([]byte, 0), fmt.Errorf(errMsg + ": %w", err)
	}
	return md, nil
}

// Main function
func RunExplainTUI() error {
	markdownFile, loadErr := loadExplainMarkdown()
	if loadErr != nil { return loadErr }

	// Render Markdown with Glamour
	glammedOutput, renderErr := glamouriseMarkdown(markdownFile) 
	if renderErr != nil {
		return fmt.Errorf("failed to glamourise the file: %w", renderErr)
	}

	p := tea.NewProgram(model{ content: glammedOutput }, tea.WithAltScreen(), tea.WithMouseCellMotion())
	// Use full size of terminal and allow mouse scrolling
	if _, err := p.Run(); err != nil {
			fmt.Println("could not run program:", err)
			os.Exit(1)
		}
		return nil
}
