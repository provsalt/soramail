package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"github.com/provsalt/soramail/internal/config"
	"github.com/provsalt/soramail/internal/tui"
	"log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	api, err := cloudflare.NewWithAPIToken(cfg.APIKey)
	if err != nil {
		panic(err)
	}

	mainMenu := tui.NewZoneMenu("Select a zone", api, nil)

	p := tea.NewProgram(mainMenu)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
