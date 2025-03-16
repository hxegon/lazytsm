package ui

import (
	"lazytsm/util"
	"unicode"

	"github.com/charmbracelet/bubbles/v2/key"
	tlist "github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	List tlist.Model
	Tmux util.Tmux
	Keys *AdditionalKeyMap
}

// changes the index of the selected item relative to the current one. Changes the filter state to FilterApplied if in filter state.
func (m *Model) slideSelectedItem(offset int) {
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

// Helper for Model.Update
func isAlphaNum(s string) bool {
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func (m Model) Init() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		switch {
		case key.Matches(msg, m.List.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.List.KeyMap.CursorDown):
			m.slideSelectedItem(1)
			return m, nil
		case key.Matches(msg, m.List.KeyMap.CursorUp):
			m.slideSelectedItem(-1)
			return m, nil
		case key.Matches(msg, m.Keys.esc):
			// Only use esc to quit in when not filtering
			if m.List.FilterState() == tlist.Unfiltered {
				return m, tea.Quit
			}
		case key.Matches(msg, m.Keys.enter):
			selected := m.List.SelectedItem().(Item).Path()
			m.Tmux.OpenOrSwitchTmuxSession(selected, selected)
		}
	}

	// Send any unhandled keys to the default list Update implementation
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// TODO: Handle case where no projects are found
func (m Model) View() string {
	return docStyle.Render(m.List.View())
}

func NewModel(items []tlist.Item) Model {
	list := tlist.New(items, tlist.NewDefaultDelegate(), 0, 0)
	list.Title = "Projects"

	list.AdditionalFullHelpKeys = KeyBindings
	list.KeyMap = *newDefaultKeyMap()

	tm, err := util.NewTmux()
	if err != nil {
		panic(err)
	}

	return Model{
		List: list,
		Tmux: tm,
		Keys: newAdditionalKeyMap(),
	}
}
