package lru_with_expirity_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	lru "github.com/cachemanager/lru_with_expirity"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAddValueWithDefaultExpiration(t *testing.T) {
	cache := lru.NewCache(1, 3)
	cache.AddWithDefaultExpiration("key", rand.Int())
}

func TestAddValueWithCustomExpiration(t *testing.T) {
	cache := lru.NewCache(1, 0)
	expectedResult := rand.Int()
	cache.AddWithCustomExpiration("key", expectedResult, 10)
	actualResult, _ := cache.Get("key")
	if actualResult != expectedResult {
		t.Fatalf("actual result: %v not we expect: %v", actualResult, expectedResult)
	}
}

func TestGetValue(t *testing.T) {
	cache := lru.NewCache(1, 1)
	expectedResult := rand.Int()
	cache.AddWithDefaultExpiration("key", expectedResult)

	if actualResult, ok := cache.Get("key"); ok {
		if actualResult != expectedResult {
			t.Fatalf("actual result: %v not we expect: %v", actualResult, expectedResult)
		}
	}
}

func TestCacheGetExpirationData(t *testing.T) {
	cache := lru.NewCache(1, 10)
	cache.AddWithDefaultExpiration("key", rand.Int())
	if cache.GetExpirationData("key").IsZero() {
		t.Fatal("default expiration time is zero")
	}
}

func TestCacheGetOldest(t *testing.T) {
	cache := lru.NewCache(3, 10)
	for i := 0; i <= 3; i++ {
		cache.AddWithDefaultExpiration(fmt.Sprintf("key%v", i), i)
	}
	if cache.GetOldest() != 1 {
		t.Fatalf("actual oldest value: %v is not we expect: %v", cache.GetOldest(), 1)
	}
}

func TestCacheDeleteValue(t *testing.T) {
	t.Run("delete existing item from cache", func(t *testing.T) {
		cache := lru.NewCache(2, 10)
		cache.AddWithDefaultExpiration("key", rand.Int())
		if !cache.Delete("key") {
			t.Fatal("failed to delete item from cache")
		}
	})
	t.Run("delete non existing item from cache", func(t *testing.T) {
		cache := lru.NewCache(2, 10)
		if cache.Delete("key") {
			t.Fatal("can delete nn existing item from cache")
		}
	})
}

func TestCachePurge(t *testing.T) {
	cache := lru.NewCache(10, 10)
	cache.AddWithDefaultExpiration(fmt.Sprintf("key%v", rand.Intn(10)), rand.Int())
	cache.Purge()
	if cache.GetCurrentSize() != 0 {
		t.Fatal("failed to purge cache")
	}
}
