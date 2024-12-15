package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"soramail/internal/request"
)

type ZonesMenu struct {
	Menu
	fetchZones tea.Cmd
	api        *cloudflare.API
	loading    bool
	spinner    spinner.Model
	errMsg     error
}

func NewZoneMenu(header string, api *cloudflare.API, parent tea.Model) *ZonesMenu {
	s := spinner.New()
	s.Spinner = spinner.Moon
	return &ZonesMenu{
		Menu: Menu{
			Header: header,
			Parent: parent,
		},
		loading:    true,
		api:        api,
		fetchZones: request.FetchZonesCmd(api),
		spinner:    s,
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
			return m, tea.Quit
		}
		for _, zone := range msg.Result {
			m.Items = append(m.Items, MenuItem{
				Name:  zone.Name,
				Model: NewDestinationMenu("Email Address", m.api, zone, m),
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

func (m *ZonesMenu) View() string {
	if m.loading {
		return fmt.Sprintf("%s Loading soramail. Press q to quit\n", m.spinner.View())
	}
	if m.errMsg != nil {
		return m.errMsg.Error()
	}
	return m.Menu.View()
}
