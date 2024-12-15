package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudflare/cloudflare-go"
	"soramail/internal/request"
	"strings"
)

type DestinationMenu struct {
	Header            string
	Items             []MenuItem
	Cursor            int
	fetchDestinations tea.Cmd
	api               *cloudflare.API
	loading           bool
	spinner           spinner.Model
	errMsg            error
	Parent            tea.Model
}

func NewDestinationMenu(header string, api *cloudflare.API, account string, parent tea.Model) *DestinationMenu {
	s := spinner.New()
	s.Spinner = spinner.Moon
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &DestinationMenu{
		Header:            header,
		loading:           true,
		api:               api,
		fetchDestinations: request.FetchDestinationCmd(api, account),
		spinner:           s,
		Parent:            parent,
	}
}

func (m *DestinationMenu) Init() tea.Cmd {
	return tea.Batch(m.fetchDestinations, m.spinner.Tick)
}

func (m *DestinationMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case request.DestinationFetchMsg:
		if msg.Err != nil {
			m.errMsg = msg.Err
		}
		m.Header = "Which email would you like to forward to?"
		for _, address := range msg.Result {
			m.Items = append(m.Items, MenuItem{
				Name:  address.Email,
				Model: nil,
			})
		}
		m.loading = false
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
	// allows the spinner to update without handling it ourselves.
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *DestinationMenu) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading emails. Press esc to go back\n", m.spinner.View())
	}
	if m.errMsg != nil {
		return m.errMsg.Error()
	}
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
