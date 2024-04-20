package project

type Item struct {
	title, desc, path string
}

func NewItem(title, desc, path string) Item {
	// TODO: Validate that the path exists?
	return Item{
		title: title,
		desc:  desc,
		path:  path,
	}
}

// list item interface
func (p Item) Title() string       { return p.title }
func (p Item) Description() string { return p.desc }
func (p Item) FilterValue() string { return p.path }

func (p Item) Path() string {
	return p.path
}
