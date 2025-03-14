package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
)

func main() {
	conf, err := ReadDefaultConfig()
	if err != nil {
		fmt.Printf("Couldn't read config file ~/.lazytsm: %v", err)
		os.Exit(1)
	}

	state := NewState(&conf)
	if err := state.LoadProjects(); err != nil {
		panic(err)
	}

	m := UIModelFromItems(state.Projects)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("lazytsm encountered an error:", err)
		os.Exit(1)
	}
}
