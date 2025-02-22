package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type lruItem struct {
	key   Key
	value interface{}
}

func (c lruCache) Set(key Key, value interface{}) bool {
	listItem, ok := c.items[key]

	if !ok {
		if c.capacity <= c.queue.Len() {
			backItem := c.queue.Back()
			delete(c.items, backItem.Value.(*lruItem).key)
			c.queue.Remove(backItem)
		}

		newListItem := c.queue.PushFront(&lruItem{key, value})

		c.items[key] = newListItem
	} else if listItem != nil {
		item := listItem.Value.(*lruItem)
		item.value = value
		c.queue.MoveToFront(listItem)
	}

	return ok
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	listItem, ok := c.items[key]

	if ok && listItem != nil {
		c.queue.MoveToFront(listItem)
		return listItem.Value.(*lruItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	for key := range c.items {
		delete(c.items, key)
	}
	c.queue.Clear() // Предположим, что у вас есть метод Clear для списка
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
