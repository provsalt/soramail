package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/provsalt/soramail/internal/tui"
	"log"
)

func main() {
	setup := tui.NewSetup()
	p := tea.NewProgram(setup)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
