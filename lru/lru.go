package lru

import (
	"container/list"
	"sync"
)

type cache struct {
	items       map[string]*list.Element
	entryList   *list.List
	lock        sync.RWMutex
	initialSize int
	size        int
}

type entry struct {
	key   string
	value interface{}
}

func NewCache(size int) *cache {
	return &cache{
		initialSize: size,
		entryList:   list.New(),
		items:       make(map[string]*list.Element, size),
	}
}

func (c *cache) Add(key string, value interface{}) (evicted bool) {
	c.lock.Lock()
	evicted = c.add(key, value)
	c.lock.Unlock()

	return evicted
}

func (c *cache) Get(key string) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.get(key)
	c.lock.Unlock()

	return value, ok
}

func (c *cache) Delete(key string) bool {
	c.lock.Lock()
	deleted := c.deleteEntry(key)
	c.lock.Unlock()

	return deleted
}

func (c *cache) GetOldest() (value interface{}) {
	c.lock.Lock()
	element := c.entryList.Back()
	value = element.Value.(*entry).value
	c.lock.Unlock()

	return value
}

func (c *cache) GetCurrentSize() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.size
}

func (c *cache) Purge() {
	c.lock.Lock()
	c.purge()
	c.lock.Unlock()
}

func (c *cache) purge() {
	c.entryList = list.New()
	c.items = make(map[string]*list.Element, c.initialSize)
	c.size = 0
}

func (c *cache) removeOldest() {
	element := c.entryList.Back()
	delete(c.items, element.Value.(*entry).key)
	c.entryList.Remove(element)
}

func (c *cache) deleteEntry(key string) bool {

	if element, ok := c.items[key]; ok && element != nil {
		c.entryList.Remove(element)
		delete(c.items, key)
		if c.size > 0 {
			c.size--
		}

		return true
	}

	return false
}

func (c *cache) get(key string) (value interface{}, ok bool) {
	if element, ok := c.items[key]; ok && element != nil {
		c.entryList.MoveToFront(element)

		return element.Value.(*entry).value, true
	}

	return nil, false
}

func (c *cache) add(key string, value interface{}) bool {
	var evicted bool

	if len(c.items) >= c.initialSize {
		c.removeOldest()
		evicted = true
		c.size--
	}

	if element, ok := c.items[key]; ok && value != nil {
		c.entryList.MoveToFront(element)

		e := element.Value.(*entry)
		e.value = value
		c.size++

		return evicted
	}

	element := c.entryList.PushFront(&entry{key: key, value: value})
	c.items[key] = element
	c.size++

	return evicted
}
