package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/provsalt/soramail/internal/config"
)

type Setup struct {
	config     config.Config
	api        *cloudflare.API
	error      error
	wizard     bool
	wizardStep int
}

type SetupInitMsg struct {
	api    *cloudflare.API
	config config.Config
	err    error
}

func NewSetup() *Setup {
	return &Setup{}
}

func (s *Setup) Init() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.ReadConfig()
		if err != nil {
			return SetupInitMsg{err: err}
		}
		if cfg.APIKey != "" {
			api, err := cloudflare.NewWithAPIToken(cfg.APIKey)
			if err != nil {
				panic(err)
			}
			return SetupInitMsg{
				api:    api,
				config: cfg,
				err:    nil,
			}
		}
		return SetupInitMsg{
			config: cfg,
		}
	}
}

func (s *Setup) processSetupKeys(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s.wizard {
		switch msg := msg.(type) {
		default:
			return s.processSetupKeys(msg)
		}
	}
	switch msg := msg.(type) {
	case SetupInitMsg:
		if msg.err != nil {
			s.error = msg.err
			return s, nil
		}
		if msg.config.APIKey == "" || msg.api == nil {
			s.wizard = true
			s.config = msg.config
			return s, nil
		}
		s.config = msg.config
		s.api = msg.api
		mainMenu := NewZoneMenu("Select a zone", s.api, nil)
		return mainMenu, mainMenu.Init()
	default:
		return s.processSetupKeys(msg)
	}
}

func (s *Setup) wizardView() string {
	return "Welcome to the setup wizard!\n"
}

func (s *Setup) View() string {
	if s.error != nil {
		return "An error occurred! " + s.error.Error()
	}
	if s.wizard {
		return s.wizardView()
	}

	// how did we end up here?
	return ""
}
