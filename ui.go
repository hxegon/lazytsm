package main

import (
	tlist "github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type UIModel struct {
	List tlist.Model
	Tmux Tmux
}

func (m UIModel) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.List.SelectedItem().(Item).Path()
			m.Tmux.OpenOrSwitchTmuxSession(selected, selected)
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// TODO: Handle case where no projects are found
func (m UIModel) View() string {
	return docStyle.Render(m.List.View())
}

func UIModelFromItems(items []Item) UIModel {
	var teaItems []tlist.Item

	// Convert to interface type, go doesn't allow for implicit conversions of pointer types
	for _, i := range items {
		teaItems = append(teaItems, tlist.Item(i))
	}

	list := tlist.New(teaItems, tlist.NewDefaultDelegate(), 0, 0)
	list.Title = "Projects"

	tm, err := NewTmux()
	if err != nil {
		panic(err)
	}

	return UIModel{
		List: list,
		Tmux: tm,
	}
}
