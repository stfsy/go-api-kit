package utils

import (
	"reflect"
	"testing"
	// Removed unused import "time"
)

type testStruct struct{}

func TestLimitedCache_StoreAndLoad(t *testing.T) {
	t1 := reflect.TypeOf(testStruct{})
	t2 := reflect.TypeOf(struct{ X int }{})
	t3 := reflect.TypeOf(struct{ Y string }{})
	t4 := reflect.TypeOf(struct{ Z float64 }{})
	cache := NewLimitedCache(2)

	// Removed redeclaration of t1, t2, t3, t4
	// t1 := reflect.TypeOf(testStruct{})
	// t2 := reflect.TypeOf(struct{ X int }{})
	// t3 := reflect.TypeOf(struct{ Y string }{})
	// t4 := reflect.TypeOf(struct{ Z float64 }{})

	cache.Store(t1, "one")
	cache.Store(t2, "two")
	cache.Store(t3, "three")

	// After storing three items in a cache of size 2, t1 should be evicted (FIFO)
	if _, ok := cache.Load(t1); ok {
		t.Error("expected t1 to be evicted after storing t3")
	}
	if v, ok := cache.Load(t2); !ok || v != "two" {
		t.Errorf("expected 'two', got %v", v)
	}
	if v, ok := cache.Load(t3); !ok || v != "three" {
		t.Errorf("expected 'three', got %v", v)
	}

	// Add fourth, should evict t2
	cache.Store(t4, "four")
	if _, ok := cache.Load(t2); ok {
		t.Error("expected t2 to be evicted after storing t4")
	}
	if v, ok := cache.Load(t3); !ok || v != "three" {
		t.Errorf("expected 'three', got %v", v)
	}
	if v, ok := cache.Load(t4); !ok || v != "four" {
		t.Errorf("expected 'four', got %v", v)
	}
}

func TestLimitedCache_UpdateMovesToEnd(t *testing.T) {
	t1 := reflect.TypeOf(testStruct{})
	t2 := reflect.TypeOf(struct{ X int }{})
	cache := NewLimitedCache(2)

	cache.Store(t1, "one")
	cache.Store(t2, "two")
	// Update t1, should move t1 to end
	cache.Store(t1, "one-updated")
	// Add a new key, should evict t2 (FIFO)
	t3 := reflect.TypeOf(struct{ Y int }{})
	cache.Store(t3, "three")
	if _, ok := cache.Load(t2); ok {
		t.Error("expected t2 to be evicted after storing t3")
	}
	if v, ok := cache.Load(t1); !ok || v != "one-updated" {
		t.Errorf("expected 'one-updated', got %v", v)
	}
	if v, ok := cache.Load(t3); !ok || v != "three" {
		t.Errorf("expected 'three', got %v", v)
	}
}
