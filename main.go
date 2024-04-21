package main

import (
	"fmt"
	"lazyproj/project"

	"github.com/charmbracelet/bubbles/list"
	tlist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
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

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var testSearchDir = "/home/hxegon/Code"

func main() {
	projs, err := project.FindGitProjects(testSearchDir)

	// TODO: Handle case where no projects are found
	if err != nil {
		fmt.Println("lazyproj has encountered an error:", err)
	}

	// Convert proj list to the right slice type
	items := make([]tlist.Item, len(projs))
	for i, p := range projs {
		items[i] = tlist.Item(p)
	}

	m := model{
		list: tlist.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	m.list.Title = "Projects"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("lazyproj encountered an error:", err)
		os.Exit(1)
	}
}
