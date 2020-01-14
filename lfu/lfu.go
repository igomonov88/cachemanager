package lfu

import (
	"container/list"
	"sync"
)

type cache struct {
	lock                sync.Mutex
	items               map[string]*list.Element
	frequencyItems      map[uint]*list.List
	startFrequency      uint
	currentMinFrequency uint
	capacity            uint
	size                uint
}

func NewCache(capacity uint) *cache {
	return &cache{
		items:               make(map[string]*list.Element),
		frequencyItems:      make(map[uint]*list.List),
		capacity:            capacity,
		startFrequency:      1,
		currentMinFrequency: 1,
	}
}

type entry struct {
	key       string
	value     interface{}
	frequency uint
}

func (c *cache) Add(key string, value interface{}) {
	c.lock.Lock()
	c.add(key, value)
	c.lock.Unlock()
}

func (c *cache) Get(key string) (value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.get(key)
}

func (c *cache) Delete(key string) {
	c.lock.Lock()
	c.deleteItem(key)
	c.lock.Unlock()
}

func (c *cache) GetCurrentSize() uint {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.currentSize()
}

func (c *cache) add(key string, value interface{}) {
	if c.size > c.capacity {
		c.evict()
	}

	if item, ok := c.items[key]; ok {
		entry := item.Value.(*entry)
		oldList := c.frequencyItems[entry.frequency]
		oldList.Remove(item)

		entry.value = value
		entry.frequency++

		if c.frequencyItems[entry.frequency] == nil {
			c.frequencyItems[entry.frequency] = list.New()
		}

		c.frequencyItems[entry.frequency].PushFront(item)
		c.size++
		return
	}

	if c.frequencyItems[c.startFrequency] == nil {
		c.frequencyItems[c.startFrequency] = list.New()
	}

	elm := c.frequencyItems[1].PushFront(&entry{key: key, value: value, frequency: 1})
	c.items[key] = elm
	c.size++

}

func (c *cache) get(key string) interface{} {
	if item, ok := c.items[key]; ok {
		entry := item.Value.(*entry)
		entry.frequency++

		if c.frequencyItems[entry.frequency] == nil {
			c.frequencyItems[entry.frequency] = list.New()
		}

		c.frequencyItems[entry.frequency].PushFront(item)
		return entry.value
	}
	return nil
}

func (c *cache) evict() {
	li := c.frequencyItems[c.currentMinFrequency]
	if li != nil && li.Len() != 0 {
		c.delete(li)
		return
	}
	c.currentMinFrequency++
	c.evict()
}

func (c *cache) delete(li *list.List) {
	var elm *list.Element

	if li != nil {
		elm = li.Back()
		ent := elm.Value.(*entry)

		if li.Len() > 1 {
			delete(c.items, ent.key)
			li.Remove(elm)
			c.size--
			return
		}

		c.frequencyItems[c.currentMinFrequency] = nil
		delete(c.items, ent.key)
		c.size--
	}
}

func (c *cache) deleteItem(key string) {
	if item, ok := c.items[key]; ok {
		ent := item.Value.(*entry)

		li := c.frequencyItems[ent.frequency]
		if li != nil {
			li.Remove(item)
		}

		delete(c.items, key)
	}
}

func (c *cache) currentSize() uint {
	return c.size
}
