package lru_test

import (
	"github.com/cachemanager/lru"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// test cases:

// delete value from cache with key which is not exist in cache

func TestLruCacheAddValue(t *testing.T) {
	cache := lru.NewCache(1)

	if ok := cache.Add("key", rand.Int()); !ok {
		t.Fatal("failed to add value to cache")
	}
}

func TestLruCacheAddValueToFullCache(t *testing.T) {
	cache := lru.NewCache(1)

	if ok := cache.Add("key", rand.Int()); ok {
		if full := cache.Add("key", rand.Int()); full {
			t.Fatalf("able to add value to cache which is over the limit. Current cache size: %v", cache.GetCurrentSize())
		}
	} else {
		t.Fatal("failed to add value to cache")
	}
}

func TestLruCacheAddValueWithExistingKey(t *testing.T) {
	cache := lru.NewCache(1)
	expectedValue := rand.Int()

	if ok := cache.Add("key", rand.Int()); ok {
		_ = cache.Add("key", expectedValue)

		if cache.GetCurrentSize() != 1 {
			t.Fatal("cache size is not we expected")
		}

		if actualValue, ok := cache.Get("key"); ok && actualValue == expectedValue {
			t.Fatalf("actual value: %v is not equal expected: %v", actualValue, expectedValue)
		}
	} else {
		t.Fatal("failed to add value to cache")
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
	cache := lru.NewCache(2)
	_ = cache.Add("key", "value")
	if deleted := cache.Delete("key"); !deleted {
		t.Fatal("delete method does not delete value from key")
	}
}

func TestLruCacheDeleteElementNotExist(t *testing.T) {
	cache := lru.NewCache(2)
	if notFound := cache.Delete("key"); !notFound && cache.GetCurrentSize() != 2 {
		t.Fatal("delete method does delete value from key")
	}
}