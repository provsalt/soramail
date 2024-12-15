package request

import (
	"context"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"soramail/internal/utils"
)

type DestinationFetchMsg struct {
	Err    error
	Result []cloudflare.EmailRoutingDestinationAddress
}

func FetchDestinationCmd(api *cloudflare.API, account string) tea.Cmd {
	return func() tea.Msg {
		return fetchDestination(api, account)
	}
}

func fetchDestination(api *cloudflare.API, account string) tea.Msg {
	if api == nil {
		return ZoneFetchedMsg{Err: errors.New("error: no cloudflare api used")}
	}
	opt := cloudflare.ListEmailRoutingAddressParameters{Verified: utils.Pointer(true)}
	emails, _, err := api.ListEmailRoutingDestinationAddresses(context.Background(), cloudflare.AccountIdentifier(account), opt)
	if err != nil {
		return DestinationFetchMsg{Err: err}
	}
	return DestinationFetchMsg{
		Result: emails,
	}
}
