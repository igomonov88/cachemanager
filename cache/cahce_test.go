package cache

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestCacheSet(t *testing.T) {
	c := New(1)
	err := c.Set("key", rand.Int())
	if err != nil {
		t.Logf("test failed due to error %v", err)
		t.Fail()
	}
}

func TestCacheSetOverTheLimit(t *testing.T) {
	c := New(0)
	err := c.Set("key", rand.Int())
	if err != ErrOverCacheLimit {
		t.Logf("test failed, got unexpected result, but expect: %v", ErrOverCacheLimit)
		t.Fail()
	}
}

func TestCacheGet(t *testing.T) {
	c := New(1)
	expectedResult := rand.Int()
	err := c.Set("key", expectedResult)
	if err != nil {
		t.Logf("test failed due to error %v", err)
		t.Fail()
	}
	actualResult, err := c.Get("key")
	if err != nil {
		t.Logf("test failed due to error %v", err)
		t.Fatal()
	}
	if expectedResult != actualResult {
		t.Logf("actual result:%v not equal expected result: %v", actualResult, expectedResult)
		t.Fail()
	}
}

func TestCacheGetNoValueWithSuchKey(t *testing.T) {
	c := New(1)
	_, err := c.Get("key")
	if err != ErrNoValueForGivenKey {
		t.Logf("expect %v error, got: %v", ErrNoValueForGivenKey, err)
	}
}

func TestCacheUpdate(t *testing.T) {
	c := New(1)
	expectedResult := rand.Int()
	_ = c.Set("key", rand.Int())
	err := c.Update("key", expectedResult)
	if err != nil {
		t.Logf("test failed due to error %v", err)
		t.Fatal()
	}
	actualResult, err := c.Get("key")
	if err != nil {
		t.Logf("test failed due to error %v", err)
		t.Fatal()
	}
	if actualResult != expectedResult {
		t.Logf("actual result:%v not equal expected result: %v", actualResult, expectedResult)
		t.Fail()
	}
}

func TestCacheUpdateNoValueForGivenKey(t *testing.T) {
	c := New(1)
	err := c.Update("key", rand.Int())
	if err != ErrNoValueForGivenKey {
		t.Logf("expect %v error, got: %v", ErrNoValueForGivenKey, err)
	}
}
