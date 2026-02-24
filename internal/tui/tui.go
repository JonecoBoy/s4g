// Package tui provides the Bubbletea post-build summary TUI for s4g.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/JonecoBoy/s4g/internal/core"
)

// ---- styles ----------------------------------------------------------------

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#059669"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DC2626"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	pageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA"))
)

// ---- model -----------------------------------------------------------------

// Model is the Bubbletea model for the build summary.
type Model struct {
	results []core.BuildResult
	done    bool
}

// New creates a new TUI model from the generator's build results.
func New(results []core.BuildResult) Model {
	return Model{results: results}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model. Pressing any key quits.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c", "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("⚡ s4g — build complete") + "\n\n")

	totalPages := 0
	totalErrors := 0

	for _, r := range m.results {
		source := fmt.Sprintf("▸ source: %s  renderer: %s", r.SourceName, r.RendererName)
		sb.WriteString(sectionStyle.Render(source) + "\n")

		for _, p := range r.Pages {
			sb.WriteString("  " + pageStyle.Render("✓") + " " + dimStyle.Render(p.Slug+".html") + "\n")
			totalPages++
		}

		for _, e := range r.Errors {
			sb.WriteString("  " + errorStyle.Render("✗ "+e.Error()) + "\n")
			totalErrors++
		}
		sb.WriteString("\n")
	}

	summary := fmt.Sprintf("%d page(s) generated", totalPages)
	if totalErrors > 0 {
		sb.WriteString(errorStyle.Render(fmt.Sprintf("⚠  %d error(s)", totalErrors)) + "  ")
	} else {
		sb.WriteString(successStyle.Render("✔  all good! "))
	}
	sb.WriteString(dimStyle.Render(summary) + "\n\n")
	sb.WriteString(dimStyle.Render("press any key to exit"))

	return sb.String()
}

// Run starts the Bubbletea program with the build summary.
// If no TTY is available (e.g. CI or piped output), it falls back to plain-text output.
func Run(results []core.BuildResult) error {
	p := tea.NewProgram(New(results))
	_, err := p.Run()
	if err != nil {
		// Bubbletea needs a terminal. Fall back to plain text in non-TTY environments.
		printPlain(results)
	}
	return nil
}

// printPlain writes a simple text summary to stdout when no TUI is available.
func printPlain(results []core.BuildResult) {
	fmt.Println("\n⚡ s4g — build complete")
	for _, r := range results {
		fmt.Printf("\nsource: %s  renderer: %s\n", r.SourceName, r.RendererName)
		for _, p := range r.Pages {
			fmt.Printf("  ✓ %s.html\n", p.Slug)
		}
		for _, e := range r.Errors {
			fmt.Printf("  ✗ %v\n", e)
		}
	}
	fmt.Println()
}
