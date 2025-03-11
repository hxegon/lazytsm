package main

import (
	"os"
)

const (
	projectsSliceOverhead = 50
)

type State struct {
	Config   *Config
	Home     string
	Tmux     *Tmux
	Projects []Item
	// UImodel?
}

func NewState(cfg *Config) State {
	tmux, err := NewTmux()

	// TODO: Proper error handling
	if err != nil {
		panic(err)
	}

	return State{
		Config:   cfg,
		Home:     os.Getenv("HOME"),
		Tmux:     &tmux,
		Projects: make([]Item, 0, len(cfg.ExtraDirs)+projectsSliceOverhead),
	}
}

func (state *State) LoadProjects() error {
	// Load extra dirs first because they're usually important
	for _, path := range state.Config.ExtraDirs {
		state.Projects = append(state.Projects, ItemFromPath(path))
	}

	gitDirs, err := FindGitProjects(state.Config.GitSearchDirs...)
	if err != nil {
		return err
	}

	state.Projects = append(state.Projects, gitDirs...)

	sessName := state.Tmux.CurrentSessionName()

	// Only include a project if it's not one currently used
	// TODO: Deduplicate projects
	lastPos := 0
	for _, p := range state.Projects {
		// TODO: Needs improvement
		if p.Description() != sessName {
			p.AbbrevPath(state.Home) // Abbreviate $HOME parts of item paths while we're at it
			state.Projects[lastPos] = p
			lastPos++
		}
	}

	state.Projects = state.Projects[:lastPos]

	return nil
}
