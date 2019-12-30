package lru_with_expirity

import (
	"container/list"
	"sync"
	"time"
)

type cache struct {
	entryList         *list.List
	items             map[string]*list.Element
	lock              sync.RWMutex
	initialSize       int
	size              int
	defaultExpiration int64
}

type entry struct {
	key        string
	value      interface{}
	expiriesAt time.Time
}

func NewCache(size int, defaultExpiration int64) *cache {
	return &cache{
		entryList:         list.New(),
		items:             make(map[string]*list.Element, size),
		lock:              sync.RWMutex{},
		initialSize:       size,
		size:              0,
		defaultExpiration: defaultExpiration,
	}
}

func (c *cache) AddWithDefaultExpiration(key string, value interface{}) (evicted bool) {
	c.lock.Lock()
	evicted = c.add(key, value, time.Duration(c.defaultExpiration)*time.Second)
	c.lock.Unlock()

	return evicted
}

func (c *cache) AddWithCustomExpiration(key string, value interface{}, expiresAt int64) (evicted bool) {
	c.lock.Lock()
	evicted = c.add(key, value, time.Duration(expiresAt)*time.Second)
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

func (c *cache) GetExpirationData(key string) time.Time {
	c.lock.RLock()
	expData := c.getExpirationData(key)
	c.lock.RUnlock()

	if expData.IsZero() {
		return time.Time{}
	}

	return expData
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
		ent := element.Value.(*entry)
		if time.Now().After(ent.expiriesAt) {
			c.deleteEntry(ent.key)
			return nil, false
		}

		c.entryList.MoveToFront(element)

		return ent.value, true
	}

	return nil, false
}

func (c *cache) add(key string, value interface{}, ttl time.Duration) bool {
	var evicted bool
	var expiresAt time.Time

	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	if len(c.items) >= c.initialSize {
		c.removeOldest()
		evicted = true
		c.size--
	}

	if element, ok := c.items[key]; ok && value != nil {
		c.entryList.MoveToFront(element)

		e := element.Value.(*entry)
		e.value = value
		e.expiriesAt = expiresAt
		c.size++

		return evicted
	}

	element := c.entryList.PushFront(&entry{key: key, value: value, expiriesAt: expiresAt})
	c.items[key] = element
	c.size++

	return evicted
}

func (c *cache) getExpirationData(key string) time.Time {
	if element, ok := c.items[key]; ok && element != nil {
		ent := element.Value.(*entry)
		return ent.expiriesAt
	}
	return time.Time{}
}
