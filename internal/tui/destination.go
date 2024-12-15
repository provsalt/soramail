package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/provsalt/soramail/internal/request"
)

type DestinationMenu struct {
	Menu
	fetchDestinations tea.Cmd
	api               *cloudflare.API
	loading           bool
	spinner           spinner.Model
	errMsg            error
	zone              cloudflare.Zone
}

func NewDestinationMenu(header string, api *cloudflare.API, zone cloudflare.Zone, parent tea.Model) *DestinationMenu {
	s := spinner.New()
	s.Spinner = spinner.Moon
	return &DestinationMenu{
		Menu: Menu{
			Header: header,
			Parent: parent,
		},
		zone:              zone,
		loading:           true,
		api:               api,
		fetchDestinations: request.FetchDestinationCmd(api, zone.Account.ID),
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
				Model: NewRandomAddressUI(m.api, m.zone, address.Email),
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
