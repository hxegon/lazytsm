package project

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func findGitDirs(root string) ([]string, error) {
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

func FindGitProjects(roots ...string) ([]Item, error) {
	projs := make([]Item, 0, len(roots)*10)

	for _, root := range roots {

		gitDirs, err := findGitDirs(root)
		if err != nil { // Good pattern?
			return []Item{}, err
		}

		for _, d := range gitDirs {
			finfo, err := os.Stat(filepath.Join(d, ".git/objects"))
			if err != nil {
				return []Item{}, fmt.Errorf("Error encountered trying to stat the .git/objects dir for %v", d)
			}

			projs = append(projs, Item{
				title: filepath.Base(d),
				desc:  d,
				// desc:  fmt.Sprintf("%v | last modified: %v", d, finfo.ModTime().Format("01/02/2006")),
				path:  d,
				mtime: finfo.ModTime(),
			})
		}
	}

	slices.SortFunc(projs, func(a, b Item) int { return b.mtime.Compare(a.mtime) })

	return projs, nil
}
