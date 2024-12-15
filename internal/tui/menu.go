package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type NavigateBackMsg struct{}

type MenuItem struct {
	Name  string
	Model tea.Model
}

type Menu struct {
	Header string
	Items  []MenuItem
	Cursor int
	Parent tea.Model
}

func (m *Menu) Init() tea.Cmd {
	return nil
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", "l":
			if len(m.Items) > 0 {
				selectedItem := m.Items[m.Cursor]
				if selectedItem.Model != nil {
					return selectedItem.Model, selectedItem.Model.Init()
				}
			}

		case "esc", "backspace", "h":
			if m.Parent != nil {
				return m.Parent, nil
			}

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		}
	}
	return m, nil
}

func (m *Menu) View() string {
	var res strings.Builder

	res.WriteString(m.Header + "\n\n")

	for i, item := range m.Items {
		cursor := " "
		if i == m.Cursor {
			cursor = ">"
		}
		res.WriteString(fmt.Sprintf("%s %s\n", cursor, item.Name))
	}
	helpStyle := lipgloss.NewStyle().Faint(true)
	res.WriteString(helpStyle.Render("\nup/down: navigate • enter: select"))
	if m.Parent != nil {
		res.WriteString(helpStyle.Render(" • esc: back"))
	}
	res.WriteString(helpStyle.Render(" • q: quit"))

	return res.String()
}
