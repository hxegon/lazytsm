package ui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tlist "github.com/charmbracelet/bubbles/v2/list"
)

type AdditionalKeyMap struct {
	enter key.Binding
	esc   key.Binding
}

func newAdditionalKeyMap() *AdditionalKeyMap {
	return &AdditionalKeyMap{
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select project and open/switch to project tmux session"),
		),
		esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Exit (filtering or whole program)"),
		),
	}
}

// We are overriding some default keys because we filter as soon as the user starts typing query text.
// By default, j k and q are used for commands (as well as others) but they won't work. So we need to override
// the default keymap to show the correct shortcuts.
func newDefaultKeyMap() *tlist.KeyMap {
	keymap := tlist.DefaultKeyMap()
	keymap.CursorDown = key.NewBinding(
		key.WithKeys("ctrl+n", "down"),
		key.WithHelp("↓/ctrl+n", "down"),
	)
	keymap.CursorUp = key.NewBinding(
		key.WithKeys("ctrl+p", "up"),
		key.WithHelp("↑/ctrl+p", "up"),
	)
	keymap.Quit = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)

	// Ensuring there are no query input keys left over from the default bindings
	keymap.GoToStart = key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "go to start"),
	)
	keymap.GoToEnd = key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "go to end"),
	)
	keymap.PrevPage = key.NewBinding(
		key.WithKeys("left", "pgup"),
		key.WithHelp("←/pgup", "prev page"),
	)
	keymap.NextPage = key.NewBinding(
		key.WithKeys("right", "pgdown"),
		key.WithHelp("→/pgdn", "next page"),
	)

	// Disable this keybinding to not show in short help. It's the binding for enter.
	// TODO: Figure out how to show enter correctly in the short help
	// I don't want to change acceptWhileFiltering b/c it could have unintended consequences. Might be overthinking it
	keymap.AcceptWhileFiltering.SetEnabled(false)

	return &keymap
}

func KeyBindings() []key.Binding {
	keyMap := newAdditionalKeyMap()
	return []key.Binding{
		keyMap.enter,
		keyMap.esc,
	}
}
