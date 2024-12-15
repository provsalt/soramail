package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/provsalt/soramail/internal/request"
	"github.com/provsalt/soramail/pkg/random"
)

type RandomAddressUI struct {
	api       *cloudflare.API
	init      tea.Cmd
	email     string
	forwarded string
	loading   bool
	zone      cloudflare.Zone
	spinner   spinner.Model
	errMsg    error
}

func NewRandomAddressUI(api *cloudflare.API, zone cloudflare.Zone, email string) *RandomAddressUI {
	s := spinner.New()
	s.Spinner = spinner.Moon
	provider := random.DefaultRandomizer{Length: 10}
	return &RandomAddressUI{
		api:     api,
		init:    request.CreateRandomAddressCmd(api, zone, email, provider),
		loading: true,
		spinner: s,
		email:   email,
		errMsg:  nil,
		zone:    zone,
	}
}

func (r *RandomAddressUI) Init() tea.Cmd {
	return r.init
}

func (r *RandomAddressUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case request.CreateRandomAddressMsg:
		if msg.Err != nil {
			r.errMsg = msg.Err
			return r, tea.Quit
		}
		r.forwarded = msg.Address
		r.loading = false
		return r, tea.Quit
	default:
		if r.loading {
			var cmd tea.Cmd
			r.spinner, cmd = r.spinner.Update(msg)
			return r, cmd
		}
	}
	return r, nil
}

func (r *RandomAddressUI) View() string {
	loadingString := fmt.Sprintf("%s Generating a random address based on domain %s for %s", r.spinner.View(), r.zone.Name, r.email)
	if r.loading {
		return loadingString
	}
	return fmt.Sprintf("%s\n ✔️ Generated successfully. Your addresse is %s\n", loadingString, r.forwarded)
}
