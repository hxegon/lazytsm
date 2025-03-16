package ui

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type uiKeyMap struct {
	navNext key.Binding
	navPrev key.Binding
	quit    key.Binding
	quit2   key.Binding
}

func newUIKeyMap() *uiKeyMap {
	return &uiKeyMap{
		navNext: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "Navigate next item"),
		),
		navPrev: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "Navigate previous item"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit program"),
		),
		quit2: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Quit program"),
		),
	}
}

func uiKeyBindings() []key.Binding {
	keyMap := newUIKeyMap()
	return []key.Binding{
		keyMap.navNext,
		keyMap.navPrev,
		keyMap.quit,
		keyMap.quit2,
	}
}
