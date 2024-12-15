package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
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
	Parent *tea.Model
}

func NewMenu(header string, items []MenuItem, parent *tea.Model) *Menu {
	return &Menu{
		Header: header,
		Items:  items,
		Parent: parent,
	}
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
				if selectedItem.Model == nil {
					return m, nil
				}
				return selectedItem.Model, nil
			}

		case "esc", "backspace", "h":
			if m.Parent != nil {
				return *m.Parent, nil
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

	res.WriteString("\nup/down: navigate • enter: select")
	if m.Parent != nil {
		res.WriteString(" • esc: back")
	}
	res.WriteString(" • q: quit")

	return res.String()
}
