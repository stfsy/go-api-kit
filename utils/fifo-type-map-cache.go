package utils

import (
	"reflect"
	"sync"
)

type LimitedCache struct {
	m      map[reflect.Type]interface{}
	keys   []reflect.Type // FIFO order
	maxLen int
	mu     sync.Mutex
}

func NewLimitedCache(maxLen int) *LimitedCache {
	return &LimitedCache{
		m:      make(map[reflect.Type]interface{}),
		maxLen: maxLen,
	}
}

func (c *LimitedCache) Store(t reflect.Type, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.m[t]; exists {
		// Remove t from keys
		for i, k := range c.keys {
			if k == t {
				c.keys = append(c.keys[:i], c.keys[i+1:]...)
				break
			}
		}
		c.m[t] = v
		// Move t to end
		c.keys = append(c.keys, t)
		return
	}
	if len(c.keys) >= c.maxLen {
		oldest := c.keys[0]
		c.keys = c.keys[1:]
		delete(c.m, oldest)
	}
	c.keys = append(c.keys, t)
	c.m[t] = v
}

func (c *LimitedCache) Load(t reflect.Type) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.m[t]
	return v, ok
}
