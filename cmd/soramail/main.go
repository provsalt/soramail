package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"soramail/internal/tui"
)

func main() {
	fmt.Println("Hello world")
	mainMenu := tui.NewMenu()

	p := tea.NewProgram(mainMenu)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
