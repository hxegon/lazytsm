package main

import (
	"lazytsm/project"
	"lazytsm/ui"
	"log/slog"
	"os"

	tlist "github.com/charmbracelet/bubbles/v2/list"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// Creates a new logger instance, need this to be more configurable later
func NewLogger(level *slog.LevelVar) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}

func main() {
	// Setup logging
	loglevel := &slog.LevelVar{}
	loglevel.Set(slog.LevelDebug) // TODO: Make log level configurable.
	// TODO: Initial log level is debug, then defaults to warn if none is set manually
	logger := NewLogger(loglevel)
	slog.SetDefault(logger)

	// Get configuration
	conf, err := ReadDefaultConfig()
	if err != nil {
		slog.Warn("Couldn't read config file ~/.lazytsm", "error", err)
	}
	slog.Debug("read configuration")

	state := NewState(&conf)
	slog.Debug("loading projects")
	if err := state.LoadProjects(); err != nil {
		slog.Error("error encountered while loading projects", "error", err)
		os.Exit(1)
	}

	// Convert projects to suitable ui items
	items := make([]tlist.Item, len(state.Projects))

	for idx, p := range state.Projects {
		items[idx] = ui.NewItem(
			project.ShortPath(p, state.Home),
			p,
		)
	}

	m := ui.NewModel(items)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if m, err := p.Run(); err != nil {
		slog.Error("lazytsm encountered an error:", "error", err)
		os.Exit(1)
	} else {
		selected := m.(ui.Model).SelectedSession
		tm := m.(ui.Model).Tmux

		if len(selected) == 0 {
			slog.Info("Nothing was selected")
			os.Exit(0)
		} else {
			tm.OpenOrSwitchTmuxSession(selected, selected)
			os.Exit(1) // Just in case
		}
	}
}
