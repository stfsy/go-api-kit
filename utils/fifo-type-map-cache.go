package utils

import (
	"reflect"
	"sync"
	"time"
)

type cacheEntry struct {
	value interface{}
	added time.Time
}

type LimitedCache struct {
	m      sync.Map
	keys   []reflect.Type
	maxLen int
	mu     sync.Mutex
}

func NewLimitedCache(maxLen int) *LimitedCache {
	return &LimitedCache{maxLen: maxLen}
}

func (c *LimitedCache) Store(t reflect.Type, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.keys) >= c.maxLen {
		oldest := c.keys[0]
		c.keys = c.keys[1:]
		c.m.Delete(oldest)
	}
	c.keys = append(c.keys, t)
	c.m.Store(t, cacheEntry{value: v, added: time.Now()})
}

func (c *LimitedCache) Load(t reflect.Type) (interface{}, bool) {
	v, ok := c.m.Load(t)
	if !ok {
		return nil, false
	}
	return v.(cacheEntry).value, true
}
