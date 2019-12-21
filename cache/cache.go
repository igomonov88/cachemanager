package cache

import "sync"

type cache struct {
	lock    sync.RWMutex
	storage map[string]interface{}
	maxSize int
	size    int
}

func New(size int) *cache {
	cache := cache{
		lock:    sync.RWMutex{},
		storage: make(map[string]interface{}, size),
		maxSize: size,
	}
	return &cache
}

func (c *cache) Set(key string, value interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.size < c.maxSize && c.maxSize != 0 {
		c.storage[key] = value
		c.size++
		return nil
	}
	return ErrOverCacheLimit
}

func (c *cache) Get(key string) (interface{}, error) {
	c.lock.RLock()
	val, ok := c.storage[key]
	c.lock.RUnlock()
	if !ok {
		return nil, ErrNoValueForGivenKey
	}
	return val, nil
}

func (c *cache) Delete(key string) error {
	c.lock.Lock()
	_, ok := c.storage[key]
	if !ok {
		c.lock.Unlock()
		return ErrNoValueForGivenKey
	}
	delete(c.storage, key)
	c.size--
	c.lock.Unlock()
	return nil
}

func (c *cache) Update(key string, value interface{}) error {
	c.lock.Lock()
	_, ok := c.storage[key]
	if !ok {
		c.lock.Unlock()
		return ErrNoValueForGivenKey
	}
	c.storage[key] = value
	c.lock.Unlock()
	return nil
}

func (c *cache) Purge() {
	c.lock.Lock()
	c.size = 0
	c.lock.Unlock()
}

func (c *cache) GetSize() int {
	c.lock.Lock()
	size := c.size
	c.lock.Unlock()
	return size
}

func (c *cache) GetMaxSize() int {
	return c.maxSize
}