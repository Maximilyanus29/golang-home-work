package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}
func (l *list) Front() *ListItem {
	return l.head
}
func (l *list) Back() *ListItem {
	return l.tail
}
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{v, nil, nil}
	current := l.Front()

	if l.Len() == 0 {
		l.head = newItem
		l.tail = newItem
	}

	newItem.Next = current
	if current != nil {
		current.Prev = newItem
	}

	l.head = newItem
	l.len++
	return newItem
}
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{v, nil, nil}
	current := l.Back()

	if l.Len() == 0 {
		l.head = newItem
		l.tail = newItem
	}

	newItem.Prev = current
	if current != nil {
		current.Next = newItem
	}

	l.tail = newItem
	l.len++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	i.Next = nil
	i.Prev = nil
	i.Value = nil
	i = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	current := l.Front()
	current.Prev = i

	if current != i {
		if i.Prev != nil {
			if i.Next != nil {
				i.Prev.Next = i.Next.Prev
			} else {
				i.Prev.Next = nil
			}
		}

		if i.Next != nil {
			if i.Prev != nil {
				i.Next.Prev = i.Prev.Next
			} else {
				i.Next.Prev = nil
			}
		}

		i.Prev = nil
		i.Next = current
		l.head = i
	}
}

func NewList() List {
	return new(list)
}
