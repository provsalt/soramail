package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/provsalt/soramail/internal/config"
	"reflect"
)

type Settings struct {
	Menu
	config      *config.Config
	configuring bool
}

func (s *Settings) Init() tea.Cmd {
	return func() tea.Msg {
		ref := reflect.ValueOf(s.config)
		iter := ref.Elem().MapRange()
		for iter.Next() {
			tea.Println(iter.Key().Interface(), iter.Value().Interface())
		}
		return nil
	}
}

func (s *Settings) updateConfig(key reflect.Value) {

}

func (s *Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if s.configuring {
				s.configuring = !s.configuring
				return s, nil
			}
			s.configuring = !s.configuring
			return s, nil
		}
	}
	return s.Menu.Update(msg)
}

func (s *Settings) View() string {
	return s.Menu.View()
}
