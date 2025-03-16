package main

import (
	"fmt"
	"lazytsm/project"
	"lazytsm/ui"
	"os"

	tlist "github.com/charmbracelet/bubbles/v2/list"
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

	// Convert projects to suitable ui items
	items := make([]tlist.Item, len(state.Projects))

	for idx, p := range state.Projects {
		items[idx] = ui.NewItem(
			project.ShortPath(p, state.Home),
			p,
		)
	}

	m := ui.NewModel(items)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if m, err := p.Run(); err != nil {
		fmt.Println("lazytsm encountered an error:", err)
		os.Exit(1)
	} else {
		selected := m.(ui.Model).SelectedSession
		tm := m.(ui.Model).Tmux

		if len(selected) == 0 {
			fmt.Println("Nothing was selected")
		}

		tm.OpenOrSwitchTmuxSession(selected, selected)
		os.Exit(0)
	}
}
