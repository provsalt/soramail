package tui

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/provsalt/soramail/internal/config"
	"reflect"
)

type Setup struct {
	config     config.Config
	textInput  textinput.Model
	api        *cloudflare.API
	error      error
	wizard     bool
	wizardStep int
	message    string
}

type SetupInitMsg struct {
	api    *cloudflare.API
	config config.Config
	err    error
}

func NewSetup() *Setup {
	ti := textinput.New()
	ti.Focus()
	return &Setup{
		textInput: ti,
	}
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

func (s *Setup) processSetupKeys(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit
		}
	}
	return nil
}

func (s *Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s.wizard {
		var cmd tea.Cmd
		switch mes := msg.(type) {
		case tea.KeyMsg:
			if mes.Type == tea.KeyEnter {
				r := reflect.ValueOf(&s.config)
				rt := r.Elem().Type()
				field := rt.Field(s.wizardStep)
				value := reflect.Indirect(r).FieldByName(field.Name)
				if !value.CanSet() {
					panic(errors.New("cannot set value to " + field.Name))
				}
				value.SetString(s.textInput.Value())
				s.textInput.Reset()
				if s.wizardStep+1 >= r.Elem().NumField() {
					var err error

					s.api, err = cloudflare.NewWithAPIToken(s.config.APIKey)

					if err != nil {
						s.error = err
						return s, tea.Quit
					}
					_, err = s.api.VerifyAPIToken(context.Background())
					if err != nil {
						s.message = "Invalid API token provided\n"
						s.wizardStep = 0
						return s, nil
					}
					err = config.SaveConfig(s.config)
					if err != nil {
						s.error = err
						return s, tea.Quit
					}

					mainMenu := NewZoneMenu("Select a zone", s.api, nil)
					return mainMenu, mainMenu.Init()
				}
				s.wizardStep++
			}
		}
		s.textInput, cmd = s.textInput.Update(msg)
		return s, tea.Batch(cmd, s.processSetupKeys(msg))
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
		return s, s.processSetupKeys(msg)
	}
}

func (s *Setup) wizardView() string {
	if s.wizardStep >= reflect.ValueOf(&s.config).Elem().NumField() {
		panic("Wizard step out of range.")
	}
	r := reflect.ValueOf(&s.config).Elem()
	rt := r.Type()
	field := rt.Field(s.wizardStep)

	return fmt.Sprintf("%s \n%sPlease enter your %s \n%s",
		"Welcome to the soramail!",
		s.message,
		field.Tag.Get("setting"),
		s.textInput.View())
}

func (s *Setup) View() string {
	if s.error != nil {
		return "An error occurred! " + s.error.Error() + "\n"
	}
	if s.wizard {
		return s.wizardView()
	}

	// how did we end up here?
	return ""
}
