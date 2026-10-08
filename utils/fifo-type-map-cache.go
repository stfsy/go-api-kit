package utils

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type LimitedCache struct {
	m      atomic.Pointer[map[reflect.Type]interface{}]
	keys   []reflect.Type // FIFO order
	maxLen int
	mu     sync.Mutex
}

func NewLimitedCache(maxLen int) *LimitedCache {
	c := &LimitedCache{
		maxLen: maxLen,
	}
	empty := make(map[reflect.Type]interface{})
	c.m.Store(&empty)
	return c
}

func (c *LimitedCache) Store(t reflect.Type, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldM := *c.m.Load()
	newM := make(map[reflect.Type]interface{}, len(oldM)+1)
	for k, val := range oldM {
		newM[k] = val
	}

	if _, exists := newM[t]; exists {
		// Remove t from keys
		for i, k := range c.keys {
			if k == t {
				c.keys = append(c.keys[:i], c.keys[i+1:]...)
				break
			}
		}
		newM[t] = v
		// Move t to end
		c.keys = append(c.keys, t)
		c.m.Store(&newM)
		return
	}

	if len(c.keys) >= c.maxLen && c.maxLen > 0 {
		oldest := c.keys[0]
		c.keys = c.keys[1:]
		delete(newM, oldest)
	}

	c.keys = append(c.keys, t)
	newM[t] = v
	c.m.Store(&newM)
}

func (c *LimitedCache) Load(t reflect.Type) (interface{}, bool) {
	mp := c.m.Load()
	if mp == nil {
		return nil, false
	}
	v, ok := (*mp)[t]
	return v, ok
}
