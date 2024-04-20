package main

import (
	"fmt"
	"lazyproj/project"

	"github.com/charmbracelet/bubbles/list"
	tlist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path/filepath"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list    tlist.Model
	selectC chan string
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
			m.selectC <- selected
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	dirs, err := project.FindGitDirs("/home/hxegon/Code")
	if err != nil {
		fmt.Println(err)
	}

	projList := make([]tlist.Item, 0, len(dirs))

	for _, path := range dirs {
		projList = append(projList, project.NewItem(filepath.Base(path), path))
	}

	selectC := make(chan string)

	m := model{
		list:    tlist.New(projList, list.NewDefaultDelegate(), 0, 0),
		selectC: selectC,
	}
	m.list.Title = "Projects"

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		selected := <-selectC
		fmt.Println("Selected project:", selected)
		os.Exit(0)
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("lazyproj encountered an error:", err)
		os.Exit(1)
	}

	// for _, p := range dirs {
	// 	fmt.Printf("%v\n", p)
	// }
}
