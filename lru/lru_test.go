package lru

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewCache(t *testing.T) {
	lru := NewCache(2)
	lru.Add("key", 12)
	lru.Add("newKey", 13)
	_, _ = lru.Get("key")
	_, _ = lru.Get("newKey")
	_, _ = lru.Get("newKey")
	t.Log(lru.entryList.Front().Value)
}

// test cases:

// get value from cache with key which is not in cache

// delete value from cache
// delete value from cache with key which is not exist in cache

func TestLruCacheAddValue(t *testing.T) {
	lru := NewCache(1)

	if ok := lru.Add("key", rand.Int()); !ok {
		t.Fatal("failed to add value to cache")
	}
}

func TestLruCacheAddValueToFullCache(t *testing.T) {
	lru := NewCache(1)

	if ok := lru.Add("key", rand.Int()); ok {
		if full := lru.Add("key", rand.Int()); full {
			t.Fatalf("able to add value to cache which is over the limit. Current cache size: %v", lru.GetCurrentSize())
		}
	} else {
		t.Fatal("failed to add value to cache")
	}
}

func TestLruCacheAddValueWithExistingKey(t *testing.T) {
	lru := NewCache(1)
	expectedValue := rand.Int()

	if ok := lru.Add("key", rand.Int()); ok {
		_ = lru.Add("key", expectedValue)

		if lru.GetCurrentSize() != 1 {
			t.Fatal("cache size is not we expected")
		}

		if actualValue, ok := lru.Get("key"); ok && actualValue == expectedValue {
			t.Fatalf("actual value: %v is not equal expected: %v", actualValue, expectedValue)
		}
	} else {
		t.Fatal("failed to add value to cache")
	}
}

func TestLruCacheGetValue(t *testing.T) {
	lru := NewCache(1)
	expectedValue := rand.Int()

	_ = lru.Add("key", expectedValue)
	actualValue, _ := lru.Get("key")

	if expectedValue != actualValue {
		t.Fatalf("actual value: %v is not equal expected: %v", actualValue, expectedValue)
	}
}
