package lru

import (
	"container/list"
	"sync"
)

type cache struct {
	items     map[string]*list.Element
	entryList *list.List
	lock      sync.RWMutex
	size      int
}

func NewCache(size int) *cache {
	return &cache{
		size:      size,
		entryList: list.New(),
		items:     make(map[string]*list.Element, size),
	}
}

func (c *cache) Add(key string, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.items) >= c.size {
		return false
	}
	if entry, ok := c.items[key]; ok && value != nil {
		entry.Value = value
		c.entryList.MoveToFront(entry)
		c.items[key] = entry

		return true
	}
	element := c.entryList.PushFront(value)
	c.items[key] = element

	return true
}

func (c *cache) Get(key string) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if entry, ok := c.items[key]; ok && entry != nil {
		c.entryList.MoveToFront(entry)
		return entry.Value, true
	}

	return nil, false
}

func (c *cache) Delete(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if entry, ok := c.items[key]; ok && entry != nil {
		c.entryList.Remove(entry)
		delete(c.items, key)

		return true
	}

	return false
}

func (c *cache) GetOldest() interface{} {
	c.lock.Lock()
	value := c.entryList.Back().Value
	c.lock.Unlock()
	return value
}

func (c *cache) GetCurrentSize() int {
	c.lock.RLock()
	size := len(c.items)
	c.lock.RUnlock()
	return size
}
