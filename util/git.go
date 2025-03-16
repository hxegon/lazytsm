package util

import (
	"fmt"
	"lazytsm/project"
	"os"
	"path/filepath"
	"time"
)

type GitDir struct {
	project.Project
}

func NewGitDir(path string) project.Project {
	return &GitDir{
		Project: project.NewDefaultProject(path),
	}
}

// Override default projecter ModTime() to check modtime for .git/objects instead
func (gd GitDir) ModTime() (time.Time, error) {
	finfo, err := os.Stat(filepath.Join(gd.Path(), ".git/objects"))
	if err != nil {
		return time.Time{}, err
	}

	return finfo.ModTime(), nil
}

func GitProjectSearch(root string) ([]project.Project, error) {
	gitDirs := make([]project.Project, 0, 30) // Give the gitDirs slice a small initial length overhead

	walkDirFn := func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while accessing path: %q, %v", path, err)
			return err
		}

		// Check to see if a directory contains a .git folder
		// and if it does, don't look inside of it for more .git folders
		if entry.IsDir() {
			ls, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, e := range ls {
				if e.Name() == ".git" {
					gitDirs = append(gitDirs, NewGitDir(path))
					return filepath.SkipDir
				}
			}
		}

		return nil
	}

	if err := filepath.WalkDir(root, walkDirFn); err != nil {
		fmt.Fprintf(os.Stderr, "Error walking path: %q, %v", root, err)
		return gitDirs, err
	}

	return gitDirs, nil
}
