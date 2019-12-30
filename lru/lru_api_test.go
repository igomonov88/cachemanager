package lru_test

import (
	"fmt"
	"github.com/cachemanager/lru"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestLruCacheAddValue(t *testing.T) {
	cache := lru.NewCache(1)
	expectedResult := rand.Int()
	cache.Add("key", expectedResult)
	actualResult, found := cache.Get("key")
	if !found {
		t.Fatal("failed to get element from the cache ")
	}
	if expectedResult != actualResult {
		t.Fatalf("actual result: %v is not we expect: %v", actualResult, expectedResult)
	}
}

func TestLruCacheAddValueEvictCache(t *testing.T) {
	cache := lru.NewCache(1)
	var isEvicted bool

	for i := 0; i <= 2; i++ {
		isEvicted = cache.Add(fmt.Sprintf("key%v", i), i)
	}

	if !isEvicted {
		t.Fatalf("element from cahce is not evicted")
	}

	if cache.GetOldest() != 2 {
		t.Fatalf("actual value:%v of cache item is not we expect:%v", cache.GetOldest(), 2)
	}
}

func TestLruCacheAddValueWithExistingKey(t *testing.T) {
	cache := lru.NewCache(1)
	expectedValue := rand.Int()
	cache.Add("key", rand.Int())

	if evicted := cache.Add("key", expectedValue); !evicted {
		t.Fatal("value with the same key is not evicted")
	}

	if cache.GetCurrentSize() != 1 {
		t.Fatal("cache size is not we expected")
	}

	if actualValue, _ := cache.Get("key"); actualValue != expectedValue {
		t.Fatalf("actual value: %v is not equal expected: %v", actualValue, expectedValue)
	}

}

func TestLruCacheGetValue(t *testing.T) {
	cache := lru.NewCache(1)
	expectedValue := rand.Int()

	_ = cache.Add("key", expectedValue)
	actualValue, _ := cache.Get("key")

	if expectedValue != actualValue {
		t.Fatalf("actual value: %v is not equal expected: %v", actualValue, expectedValue)
	}
}

func TestLruCacheGetNotExistValue(t *testing.T) {
	cache := lru.NewCache(1)
	if value, exist := cache.Get("key"); exist {
		t.Fatalf("get method returns the value: %v from key which does not exist", value)
	}
}

func TestLruCacheDeleteMethod(t *testing.T) {
	cache := lru.NewCache(1)
	_ = cache.Add("key", rand.Int())
	if deleted := cache.Delete("key"); !deleted {
		t.Fatal("delete method does not delete value from key")
	}
}

func TestLruCacheDeleteElementNotExist(t *testing.T) {
	cache := lru.NewCache(1)

	if deleted := cache.Delete("key"); deleted {
		t.Fatal("delete method does delete value from key")
	}

}

func TestLruCachePurge(t *testing.T) {
	cache := lru.NewCache(2)
	cache.Add("key", "value")
	cache.Purge()

	if cache.GetCurrentSize() != 0 {
		t.Fatal("current size is not equal what we expect")
	}

	if value, ok := cache.Get("key"); ok && value != nil {
		t.Fatalf("cache is not purged")
	}
}

func TestLruCacheGetOldest(t *testing.T) {
	cache := lru.NewCache(4)
	for i := 0; i <= 3; i++ {
		cache.Add(fmt.Sprintf("key%v", i), i)
	}
	if cache.GetOldest() != 0 {
		t.Fatalf("oldest value:%v is not we expect:%v ", cache.GetOldest(), 0)
	}
}
