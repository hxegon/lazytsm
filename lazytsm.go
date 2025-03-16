package main

import (
	// "fmt"
	"lazytsm/project"
	"lazytsm/util"
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
		state.addProject(project.NewDefaultProject(path), sessName)
	}

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
		// temporarily disabling this error checking until #17 is fixed
		// amtime, errA := a.ModTime()
		// bmtime, errB := b.ModTime()
		amtime, _ := a.ModTime()
		bmtime, _ := b.ModTime()

		// if errA != nil || errB != nil {
		// 	err := fmt.Errorf("Unexpected errors when sorting projects by modification time: a: %v, b: %v", errA, errB)
		// 	panic(err)
		// }

		return bmtime.Compare(amtime)
	})

	return nil
}
