package ns

type NamespaceItem struct {
	Name string
}

func (i NamespaceItem) Title() string       { return i.Name }
func (i NamespaceItem) FilterValue() string { return i.Name }
