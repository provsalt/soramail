package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudflare/cloudflare-go"
	"soramail/internal/request"
	"strings"
)

type ZonesMenu struct {
	Header     string
	Items      []MenuItem
	Cursor     int
	fetchZones tea.Cmd
	api        *cloudflare.API
	loading    bool
	spinner    spinner.Model
	errMsg     error
	Parent     *tea.Model
}

func NewZoneMenu(header string, api *cloudflare.API, parent *tea.Model) *ZonesMenu {
	s := spinner.New()
	s.Spinner = spinner.Moon
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &ZonesMenu{
		Header:     header,
		loading:    true,
		api:        api,
		fetchZones: request.FetchZonesCmd(api),
		spinner:    s,
		Parent:     parent,
	}
}

func (m *ZonesMenu) Init() tea.Cmd {
	return tea.Batch(m.fetchZones, m.spinner.Tick)
}

func (m *ZonesMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case request.ZoneFetchedMsg:
		if msg.Err != nil {
			m.errMsg = msg.Err
		}
		for _, zone := range msg.Result {
			m.Items = append(m.Items, MenuItem{Name: zone.Name})
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
	// allows the spinner to update without handling it ourselves.
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *ZonesMenu) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading soramail. Press q to quit\n", m.spinner.View())
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
