package ui

import (
	tlist "github.com/charmbracelet/bubbles/v2/list"
	"lazytsm/project"
	"path/filepath"
)

type Item struct {
	project.Project
	desc string
}

func NewItem(desc string, project project.Project) tlist.Item {
	// TODO: Validate that the path exists?
	return Item{
		Project: project,
		desc:    desc,
	}
}

// list item interface
func (i Item) Title() string       { return filepath.Base(i.Path()) }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.Path() }
