package project

import (
	tlist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type UIModel struct {
	list tlist.Model
}

func (m UIModel) Init() tea.Cmd {
	return nil // do nothing for now
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(Item).Path()
			t, _ := NewTmux() // TODO handle error
			t.OpenOrSwitchTmuxSession(selected, selected)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// TODO: Handle case where no projects are found
func (m UIModel) View() string {
	return docStyle.Render(m.list.View())
}
