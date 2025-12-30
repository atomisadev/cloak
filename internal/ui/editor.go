package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	NeonCyan   = "#00FFFF"
	NeonPink   = "#FF00FF"
	DeepVoid   = "#0A0A0A"
	GhostWhite = "#E0E0E0"
	AlertRed   = "#FF3333"
	MutedGray  = "#444444"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(MutedGray)).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(NeonCyan)).
			Bold(true).
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color(MutedGray))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DeepVoid)).
			Background(lipgloss.Color(NeonCyan)).
			Bold(true)

	dimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(MutedGray))

	inputPopupStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color(NeonPink)).
			Padding(1, 2).
			Align(lipgloss.Center)
)

type AppState int

const (
	StateBrowsing AppState = iota
	StateEditingValue
	StateAddingKey
	StateConfirmDelete
)

type KeyValue struct {
	Key   string
	Value string
}

type Model struct {
	State      AppState
	Table      table.Model
	Input      textinput.Model
	Secrets    []KeyValue
	ShowValues bool
	ColFocus   int
	ToSave     map[string]string
	Quitting   bool
}

func InitialModel(secrets map[string]string) Model {
	var data []KeyValue
	for k, v := range secrets {
		data = append(data, KeyValue{Key: k, Value: v})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].Key < data[j].Key
	})

	columns := []table.Column{
		{Title: "KEY", Width: 20},
		{Title: "VALUE (Masked)", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = headerStyle
	s.Selected = selectedStyle
	s.Cell = lipgloss.NewStyle().Padding(0, 1)
	t.SetStyles(s)

	ti := textinput.New()
	ti.CharLimit = 2048
	ti.Width = 50
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(NeonPink))

	m := Model{
		State:      StateBrowsing,
		Table:      t,
		Input:      ti,
		Secrets:    data,
		ShowValues: false,
		ColFocus:   1,
	}
	m.updateTableRows()
	return m
}

func (m *Model) updateTableRows() {
	var rows []table.Row
	for _, s := range m.Secrets {
		valDisplay := "••••••••••••"
		if m.ShowValues {
			valDisplay = s.Value
		}
		rows = append(rows, table.Row{s.Key, valDisplay})
	}
	m.Table.SetRows(rows)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Table.SetWidth(msg.Width - 4)
		m.Table.SetHeight(msg.Height - 10)

	case tea.KeyMsg:
		switch m.State {
		case StateBrowsing:
			switch msg.String() {
			case "q", "ctrl+c":
				m.Quitting = true
				m.ToSave = nil
				return m, tea.Quit
			case "ctrl+s":
				m.Quitting = true
				m.ToSave = m.exportMap()
				return m, tea.Quit
			case "v", " ":
				m.ShowValues = !m.ShowValues
				m.updateTableRows()
			case "d", "backspace":
				if len(m.Secrets) > 0 {
					m.State = StateConfirmDelete
				}
			case "a":
				m.State = StateAddingKey
				m.Input.Placeholder = "NEW_KEY_NAME"
				m.Input.SetValue("")
				m.Input.Focus()
				return m, textinput.Blink
			case "enter":
				if len(m.Secrets) > 0 {
					m.State = StateEditingValue
					row := m.Table.Cursor()
					m.Input.Placeholder = "Value..."
					m.Input.SetValue(m.Secrets[row].Value)
					m.Input.Focus()
					return m, textinput.Blink
				}
			}

		case StateEditingValue:
			switch msg.String() {
			case "enter":
				row := m.Table.Cursor()
				if row >= 0 && row < len(m.Secrets) {
					m.Secrets[row].Value = m.Input.Value()
					m.updateTableRows()
				}
				m.State = StateBrowsing
				m.Input.Blur()
			case "esc":
				m.State = StateBrowsing
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			return m, cmd

		case StateAddingKey:
			switch msg.String() {
			case "enter":
				newKey := m.Input.Value()
				if newKey != "" {
					exists := false
					for _, s := range m.Secrets {
						if s.Key == newKey {
							exists = true
							break
						}
					}
					if !exists {
						newItem := KeyValue{Key: newKey, Value: ""}
						m.Secrets = append(m.Secrets, newItem)
						sort.Slice(m.Secrets, func(i, j int) bool {
							return m.Secrets[i].Key < m.Secrets[j].Key
						})
						m.updateTableRows()

					}
				}
				m.State = StateBrowsing
				m.Input.Blur()
			case "esc":
				m.State = StateBrowsing
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			return m, cmd

		case StateConfirmDelete:
			switch msg.String() {
			case "y", "enter":
				row := m.Table.Cursor()
				if row >= 0 && row < len(m.Secrets) {
					m.Secrets = append(m.Secrets[:row], m.Secrets[row+1:]...)
					m.Table.SetCursor(0)
					m.updateTableRows()
				}
				m.State = StateBrowsing
			case "n", "esc":
				m.State = StateBrowsing
			}
		}
	}

	if m.State == StateBrowsing {
		m.Table, cmd = m.Table.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	banner := lipgloss.NewStyle().Foreground(lipgloss.Color(NeonPink)).Bold(true).Render("CLOAK // VAULT EDITOR")

	var status string
	switch m.State {
	case StateBrowsing:
		status = fmt.Sprintf("ROWS: %d • [a] ADD • [d] DELETE • [v] TOGGLE VISIBILITY • [enter] EDIT • [ctrl+s] SAVE & QUIT", len(m.Secrets))
	case StateEditingValue:
		status = "EDITING VALUE • [enter] CONFIRM • [esc] CANCEL"
	case StateAddingKey:
		status = "NEW KEY NAME • [enter] CONFIRM • [esc] CANCEL"
	case StateConfirmDelete:
		status = lipgloss.NewStyle().Foreground(lipgloss.Color(AlertRed)).Render("DELETE SELECTED SECRET? (y/n)")
	}
	status = dimmedStyle.Render(status)

	var content string
	if m.State == StateEditingValue || m.State == StateAddingKey {
		label := "VALUE"
		if m.State == StateAddingKey {
			label = "NEW KEY"
		}

		inputView := lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color(NeonCyan)).Render(label),
			m.Input.View(),
		)
		content = inputPopupStyle.Render(inputView)
	} else {
		content = baseStyle.Render(m.Table.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		banner,
		"\n",
		content,
		"\n",
		status,
	)
}

func (m Model) exportMap() map[string]string {
	out := make(map[string]string)
	for _, s := range m.Secrets {
		out[s.Key] = s.Value
	}
	return out
}
