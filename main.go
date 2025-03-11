package main

import (
	"fmt"
	"lazytsm/config"
	"lazytsm/project"
	"os"

	tlist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

func main() {
	conf, err := config.ReadDefaultConfig()
	if err != nil {
		fmt.Printf("Couldn't read config file ~/.lazytsm: %v", err)
		os.Exit(1)
	}

	projs := make([]project.Item, 0, len(conf.ExtraDirs)+50)

	// TODO: DEDUPLICATE DIRS
	// Put in extra dirs first as they are probably important
	for _, path := range conf.ExtraDirs {
		projs = append(projs, project.ItemFromPath(path))
	}

	gDirs, err := project.FindGitProjects(conf.GitSearchDirs...)
	if err != nil {
		fmt.Println("lazytsm has encountered an error:", err)
		os.Exit(1)
	}

	projs = append(projs, gDirs...)

	// Remove the current session from the project list
	tm, err := project.NewTmux()
	if err != nil {
		panic(err)
	}

	currentName := tm.CurrentSessionName()
	home := os.Getenv("HOME")

	// Convert proj list to the right slice type
	var items []tlist.Item
	for _, p := range projs {
		// Don't add a project for any current sessions
		if p.Description() != currentName {
			p.AbbrevPath(home)
			items = append(items, tlist.Item(p))
		}
	}

	m := UIModel{
		list: tlist.New(items, tlist.NewDefaultDelegate(), 0, 0),
	}

	m.list.Title = "Projects"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("lazytsm encountered an error:", err)
		os.Exit(1)
	}
}
