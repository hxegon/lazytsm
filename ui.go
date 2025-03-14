package main

import (
	tlist "github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"unicode"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type UIModel struct {
	List tlist.Model
	Tmux Tmux
}

// changes the index of the selected item relative to the current one. Changes the filter state to FilterApplied if in filter state.
func (m *UIModel) slideSelectedItem(offset int) {
	// Need to change from Filtering -> FilterApplied, no candidate can be selected in Filtering
	if m.List.FilterState() == tlist.Filtering {
		m.List.SetFilterState(tlist.FilterApplied)
	}

	nextIdx := m.List.Index() + offset

	// Guard against selecting an out of bounds index
	if nextIdx > 0 {
		nItems := len(m.List.Items())
		// use mod to wrap around if selecting next item at end of list
		nextIdx = nextIdx % nItems
	} else {
		nextIdx = max(0, nextIdx)
	}

	m.List.Select(nextIdx)
}

// Helper for UIModel.Update
func isAlphaNum(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

func (m UIModel) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		s := msg.String()

		// Start filtering if we start receiving alphanumeric chars in a non-filtering state
		if m.List.FilterState() != tlist.Filtering && len(s) == 1 && isAlphaNum(s) {
			m.List.SetFilterText(s) // pass the entered char as the first letter of the query
			m.List.SetFilterState(tlist.Filtering)
			return m, cmd
		}

		switch s {
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+n":
			m.slideSelectedItem(1)
			return m, nil
		case "ctrl+p":
			m.slideSelectedItem(-1)
			return m, nil
		case "esc":
			// Only use esc to quit in when not filtering
			if m.List.FilterState() == tlist.Unfiltered {
				return m, tea.Quit
			}
		case "enter":
			selected := m.List.SelectedItem().(Item).Path()
			m.Tmux.OpenOrSwitchTmuxSession(selected, selected)
		}
	}

	// Send any unhandled keys to the default list Update implementation
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
