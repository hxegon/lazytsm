package project

import (
	"os"
	"strings"
	"time"
)

type Project interface {
	Path() string
	ModTime() (time.Time, error)
}

func ShortPath(p Project, home string) string {
	return strings.Replace(p.Path(), home, "~", 1)
}

// Noting a method for default implementations of interface methods and overriding them

type DefaultProject struct {
	path string
}

func (p DefaultProject) Path() string { return p.path }

func (p DefaultProject) ModTime() (time.Time, error) {
	finfo, err := os.Stat(p.Path())
	if err != nil {
		return time.Time{}, err
	}

	return finfo.ModTime(), nil
}

func NewDefaultProject(path string) Project {
	return &DefaultProject{
		path: path,
	}
}
