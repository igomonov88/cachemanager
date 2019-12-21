package lru

import (
	"container/list"
	"sync"
	"time"
)

type lruCache struct {
	lock      sync.RWMutex
	entryList *list.List
	items     map[interface{}]*list.Element
	size      int
	len       int
}

type entry struct {
	key     interface{}
	value   interface{}
	expires time.Time
}

func NewLRU(size int) *lruCache {
	return &lruCache{
		size:      size,
		entryList: list.New(),
		items:     make(map[interface{}]*list.Element, size),
	}
}

func (l *lruCache) AddValueWithExpiresInSeconds(key, value interface{}, expiresAtSec int64) {
	l.lock.Lock()
	l.add(key, value, time.Duration(expiresAtSec)*time.Second)
	l.lock.Unlock()
}

func (l *lruCache) add(key, value interface{}, ttl time.Duration) {
	var expires time.Time
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}

	// check for existing item
	if ent, ok := l.items[key]; ok {
		l.entryList.MoveToFront(ent)
		e := ent.Value.(*entry)
		e.expires = expires
		e.value = value
		l.len++
		return
	}

	// add new entry
	ent := entry{
		key:     key,
		value:   value,
		expires: expires,
	}
	entry := l.entryList.PushFront(ent)
	l.items[key] = entry
	l.len++

	if l.entryList.Len() > l.len {
		l.removeItem(l.entryList.Back())
	}
}

func (l *lruCache) Get(key interface{}) (value interface{}, ok bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.get(key)
}

func (l *lruCache) get(key interface{}) (value interface{}, ok bool) {
	if ent, ok := l.items[key]; ok {
		e := ent.Value.(*entry)

		if !e.expires.IsZero() && time.Now().After(e.expires) {
			l.removeItem(ent)
			return nil, false
		}

		l.entryList.MoveToFront(ent)
		return e.value, false
	}
	return nil, false
}

func (l *lruCache) removeItem(element *list.Element) {
	l.entryList.Remove(element)
	kv := element.Value.(*entry)
	delete(l.items, kv.key)
}
