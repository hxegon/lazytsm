package project

type Item struct {
	title, desc string
}

func NewItem(title, desc string) Item {
	return Item{
		title: title,
		desc:  desc,
	}
}

// list item interface
func (p Item) Title() string       { return p.title }
func (p Item) Description() string { return p.desc }
func (p Item) FilterValue() string { return p.title }

func (p Item) Path() string {
	return p.desc
}
