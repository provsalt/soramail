package request

import (
	"context"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
)

type ZoneFetchedMsg struct {
	Err    error
	Result []cloudflare.Zone
}

func FetchZonesCmd(api *cloudflare.API) tea.Cmd {
	return func() tea.Msg {
		return fetchZones(api)
	}
}

func fetchZones(api *cloudflare.API) tea.Msg {
	if api == nil {
		return ZoneFetchedMsg{Err: errors.New("error: no cloudflare api used")}
	}
	zones, err := api.ListZones(context.Background())
	if err != nil {
		return ZoneFetchedMsg{Err: err}
	}
	return ZoneFetchedMsg{
		Err:    nil,
		Result: zones,
	}
}
