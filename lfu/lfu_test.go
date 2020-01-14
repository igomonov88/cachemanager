package lfu

import "testing"

func TestLFU(t *testing.T) {
	cache := NewCache(5)

	cache.Add("key", 10)
	cache.Add("key2", 11)
	cache.Add("key2", 12)
	cache.Add("key4", 9)
	cache.Add("key5", 13)
	cache.Add("key6", 15)

	//actValue := cache.Get("key")
	//t.Log(actValue)

	for i, v := range cache.frequencyItems {
		t.Log(i)
		t.Log(v.Len())
	}
}
