package main

import (
	"fmt"
	"lazytsm/config"
	"lazytsm/project"
	"os"

	tlist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list tlist.Model
}

func (m model) Init() tea.Cmd {
	return nil // do nothing for now
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(project.Item).Path()
			t, _ := project.NewTmux() // TODO handle error
			t.OpenOrSwitchTmuxSession(selected, selected)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// TODO: Handle case where no projects are found
func (m model) View() string {
	return docStyle.Render(m.list.View())
}

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

	m := model{
		list: tlist.New(items, tlist.NewDefaultDelegate(), 0, 0),
	}

	m.list.Title = "Projects"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("lazytsm encountered an error:", err)
		os.Exit(1)
	}
}
