package tui

import (
	"fmt"
	"helmwatch/internal/config"
	"helmwatch/internal/diff"
	"helmwatch/internal/helm"
	"helmwatch/internal/msg"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle       = lipgloss.NewStyle().Bold(true)
	searchMatchStyle = lipgloss.NewStyle().Background(lipgloss.Color("226")).Foreground(lipgloss.Color("0"))
)

func New(config config.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Focus()

	vp := viewport.New(0, 0)

	return Model{
		input:    ti,
		viewport: vp,
		config:   config,
	}
}

type Model struct {
	viewport viewport.Model
	input    textinput.Model

	diff string

	searchMode bool
	query      string

	config config.Config

	previousDir string
}

func (m Model) Init() tea.Cmd {
	return m.renderChartAndShowDiff()
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {

	case msg.Render:
		m.diff = message.Diff
		if message.Directory != nil {
			m.previousDir = *message.Directory
		}
		m.viewport.SetContent(highlightSearchResults(m.diff, m.query))
		return m, nil

	case msg.FileChanged:
		return m, m.renderChartAndShowDiff()

	case tea.WindowSizeMsg:
		m.viewport.Width = message.Width
		m.viewport.Height = message.Height
		return m, nil

	case tea.KeyMsg:
		if m.searchMode {
			switch message.String() {
			case "enter":
				m.query = m.input.Value()
				m.viewport.SetContent(highlightSearchResults(m.diff, m.query))
				m.searchMode = false
				return m, nil
			case "esc":
				m.searchMode = false
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(message)
			return m, cmd
		}

		switch message.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			return m, m.renderChartAndShowDiff()
		case "/":
			m.searchMode = true
			m.input.SetValue("")
			m.input.Focus()
			return m, nil
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(message)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	header := titleStyle.Render("Helm Diff (q quit | r rerender | / search)")
	if m.searchMode {
		return fmt.Sprintf("%s\n/%s\n\n%s", header, m.input.View(), m.viewport.View())
	}
	return fmt.Sprintf("%s\n\n%s", header, m.viewport.View())
}

func (m *Model) renderChartAndShowDiff() tea.Cmd {
	return func() tea.Msg {
		dir, err := helm.Template(helm.TemplateOptions{
			Chart:      m.config.Chart,
			Version:    m.config.Version,
			ValuesFile: m.config.ValuesFile,
			Exclusions: m.config.Exclusions,
		})
		if err != nil {
			return msg.NewRender(fmt.Sprintf("failed to render: %s", err), nil)
		}

		if m.previousDir == "" {
			return msg.NewRender("finished initial render", &dir)
		}

		difference := diff.Dirs(m.previousDir, dir)
		if difference == "" {
			return msg.NewRender("no changes detected", &dir)
		}

		return msg.NewRender(difference, &dir)
	}
}

func highlightSearchResults(text, query string) string {
	if query == "" {
		return text
	}

	return strings.ReplaceAll(text, query, searchMatchStyle.Render(query))
}
