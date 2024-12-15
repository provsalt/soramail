package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudflare/cloudflare-go"
	"soramail/internal/request"
)

type DestinationMenu struct {
	Menu
	fetchDestinations tea.Cmd
	api               *cloudflare.API
	loading           bool
	spinner           spinner.Model
	errMsg            error
}

func NewDestinationMenu(header string, api *cloudflare.API, account string, parent tea.Model) *DestinationMenu {
	s := spinner.New()
	s.Spinner = spinner.Moon
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &DestinationMenu{
		Menu: Menu{
			Header: header,
			Parent: parent,
		},
		loading:           true,
		api:               api,
		fetchDestinations: request.FetchDestinationCmd(api, account),
		spinner:           s,
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
		return m, nil
	default:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m.Menu.Update(msg)
	}
}

func (m *DestinationMenu) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading emails. Press esc to go back\n", m.spinner.View())
	}
	if m.errMsg != nil {
		return m.errMsg.Error()
	}
	return m.Menu.View()
}
