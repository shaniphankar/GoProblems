package Problem1

type Item struct {
	value string
	priority int
	index int
}

func NewItem(value string, priority int, index int) *Item {
	p := new(Item)
	p.value=value
	p.priority=priority
	p.index=index
	return p
}

func (item Item) GetValue() string {
	return item.value;
}

func (item Item) GetPriority() int {
	return item.priority;
}
