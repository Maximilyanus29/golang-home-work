package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
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
	if current != i {
		// Удаляем элемент i из текущего положения
		if i.Prev != nil {
			i.Prev.Next = i.Next
		}
		if i.Next != nil {
			i.Next.Prev = i.Prev
		}

		// Перемещаем элемент i в начало списка
		i.Prev = nil
		i.Next = current
		if current != nil {
			current.Prev = i
		}
		l.head = i
	}
}

func (l *list) Clear() {
	for i := l.Front(); i != nil; i = i.Next {
		i.Prev = nil
		i.Next = nil
		i = nil
	}
}

func NewList() List {
	return new(list)
}
