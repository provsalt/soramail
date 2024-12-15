package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"soramail/internal/config"
	"soramail/internal/tui"
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

	mainMenu := tui.NewZoneMenu("Zones", api, nil)

	p := tea.NewProgram(mainMenu)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
