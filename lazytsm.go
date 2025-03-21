package main

import (
	"lazytsm/project"
	"lazytsm/util"
	"log/slog"
	"os"
	"slices"
)

const (
	projectsSliceOverhead = 50
)

type State struct {
	Config   *Config
	Home     string
	Tmux     *util.Tmux
	Projects []project.Project
}

func NewState(cfg *Config) State {
	tmux, err := util.NewTmux()

	// TODO: Proper error handling
	if err != nil {
		slog.Error("fatal error occurred while initializing tmux")
		panic(err)
	}

	return State{
		Config:   cfg,
		Home:     os.Getenv("HOME"),
		Tmux:     &tmux,
		Projects: make([]project.Project, 0, len(cfg.ExtraDirs)+projectsSliceOverhead),
	}
}

// TODO: manually passing sessName because I don't want to re-call CurrentSessionName
// every time we add a project, but this is ick.
func (state *State) addProject(p project.Project, sessName string) {
	if p.Path() != sessName {
		state.Projects = append(state.Projects, p)
	}
}

func (state *State) LoadProjects() error {
	// Load extra dirs first because they're usually important
	sessName := state.Tmux.CurrentSessionName()

	for _, path := range state.Config.ExtraDirs {
		if _, err := os.Stat(path); err != nil {
			// Only log the error if the dir doesn't exist, but return if it's some unknown one.
			if os.IsNotExist(err) {
				slog.Warn("path configured in ExtraDirs could not be found", "path", path)
			} else {
				return err
			}
		} else { // Only add the project if we can os.Stat it
			state.addProject(project.NewDefaultProject(path), sessName)
		}
	}

	// Don't need to stat the git dirs because they are dynamically searched for at run time, so no risk of them not existing
	for _, gitSearchDir := range state.Config.GitSearchDirs {
		newDirs, err := util.GitProjectSearch(gitSearchDir)
		if err != nil {
			return err
		}

		for _, newDir := range newDirs {
			state.addProject(newDir, sessName)
		}
	}

	// Sort projects into reverse chronological order (by last modification time)
	slices.SortFunc(state.Projects, func(a, b project.Project) int {
		amtime, errA := a.ModTime()
		bmtime, errB := b.ModTime()

		if errA != nil || errB != nil {
			slog.Error("fatal error occurred when trying to compare project modification times", "project A error", errA, "project B error", errB)
			os.Exit(1)
		}

		return bmtime.Compare(amtime)
	})

	return nil
}
