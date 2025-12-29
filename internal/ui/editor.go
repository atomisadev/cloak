package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorNeonBlue   = "#00FFFF"
	ColorDeepPurple = "#2E004B"
	ColorDarkGray   = "#1A1A1A"
	ColorWhite      = "#FFFFFF"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(ColorNeonBlue))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorNeonBlue)).
			Bold(true).
			Align(lipgloss.Center)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

type AppState int

const (
	StateBrowsing AppState = iota
	StateEditing
)

type KeyValue struct {
	Key   string
	Value string
}

type Model struct {
	State     AppState
	Table     table.Model
	Input     textinput.Model
	Secrets   []KeyValue
	ActiveRow int
}

func InitialModel(secrets map[string]string) Model {
	var rows []table.Row
	var data []KeyValue

	for k, v := range secrets {
		rows = append(rows, table.Row{k, "*****"})
		data = append(data, KeyValue{Key: k, Value: v})
	}

	columns := []table.Column{
		{Title: "KEY", Width: 20},
		{Title: "VALUE", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ColorNeonBlue)).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(ColorNeonBlue)).
		Bold(true)
	t.SetStyles(s)

	ti := textinput.New()
	ti.Placeholder = "Enter new secret value..."
	ti.CharLimit = 1024
	ti.Width = 40
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorNeonBlue))

	return Model{
		State:   StateBrowsing,
		Table:   t,
		Input:   ti,
		Secrets: data,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Table.SetWidth(msg.Width - 10)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.State == StateBrowsing {
				return m, tea.Quit
			}

		case "enter":
			if m.State == StateBrowsing {
				m.State = StateEditing
				m.ActiveRow = m.Table.Cursor()

				if m.ActiveRow >= 0 && m.ActiveRow < len(m.Secrets) {
					currentVal := m.Secrets[m.ActiveRow].Value
					m.Input.SetValue(currentVal)
					m.Input.Focus()
					return m, textinput.Blink
				}
				return m, nil
			} else {
				newValue := m.Input.Value()
				if m.ActiveRow >= 0 && m.ActiveRow < len(m.Secrets) {
					m.Secrets[m.ActiveRow].Value = newValue

					newRows := m.Table.Rows()
					newRows[m.ActiveRow][1] = "*****"
					m.Table.SetRows(newRows)
				}

				m.State = StateBrowsing
				m.Input.Blur()
				return m, nil
			}

		case "esc":
			if m.State == StateEditing {
				m.State = StateBrowsing
				m.Input.Blur()
				return m, nil
			}
		}
	}

	if m.State == StateBrowsing {
		m.Table, cmd = m.Table.Update(msg)
	} else {
		m.Input, cmd = m.Input.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	banner := `
   ______   __    ____    ___    __ __
  / ____/  / /   / __ \  /   |  / //_/
 / /      / /   / / / / / /| | / ,<
/ /___   / /___/ /_/ / / ___ |/ /| |
\____/  /_____/\____/ /_/  |_/_/ |_|
`

	bannerRendered := headerStyle.Render(banner)

	if m.State == StateEditing {
		inputView := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorNeonBlue)).
			Padding(1, 2).
			Render(
				fmt.Sprintf("EDITING: %s\n\n%s",
					m.Secrets[m.ActiveRow].Key,
					m.Input.View(),
				),
			)

		return lipgloss.JoinVertical(lipgloss.Center, bannerRendered, inputView)
	}

	tableView := baseStyle.Render(m.Table.View())

	help := helpStyle.Render("↑/↓: Navigate • Enter: Edit • q: Quit")

	return lipgloss.JoinVertical(lipgloss.Center, bannerRendered, tableView, help)
}
