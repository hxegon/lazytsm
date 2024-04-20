package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindGitDirs(root string) ([]string, error) {
	gitDirs := make([]string, 0, 30) // initialize a slice with an initial cap of 30

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
					gitDirs = append(gitDirs, path)
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
