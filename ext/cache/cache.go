package opa

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

var mCache *MemoryCache

//Local local cache
type Local interface {
	Get(key [32]byte) (interface{}, bool)
	Set(key [32]byte, val interface{})
	Delete(key [32]byte)
}

type cacheValue struct {
	value   interface{}
	expired time.Time
}

// MemoryCache is an implemtation of Cache that stores responses in an in-memory map.
type MemoryCache struct {
	mu          sync.RWMutex
	items       map[[32]byte]cacheValue
	expDuration time.Duration
}

// Get returns the []byte representation of the response and true if present, false if not
func (c *MemoryCache) Get(key [32]byte) (interface{}, bool) {
	c.mu.RLock()
	resp, ok := c.items[key]
	c.mu.RUnlock()

	if !resp.expired.IsZero() && time.Now().After(resp.expired) {
		c.Delete(key)
		return nil, false
	}

	return resp.value, ok
}

// Set saves response resp to the cache with key
func (c *MemoryCache) Set(key [32]byte, val interface{}) {
	c.mu.Lock()
	v := cacheValue{value: val}
	if c.expDuration > 0 {
		v.expired = time.Now().Add(c.expDuration)
	}
	c.items[key] = v
	c.mu.Unlock()
}

// Delete removes key from the cache
func (c *MemoryCache) Delete(key [32]byte) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// NewMemoryCache returns a new Cache that will store items in an in-memory map
func NewMemoryCache(exp time.Duration) *MemoryCache {
	c := &MemoryCache{
		items:       map[[32]byte]cacheValue{},
		expDuration: exp,
	}
	return c
}

//LRU lru cache
type LRU struct {
	cache *lru.Cache
}

//NewLRU create lru cache instance
func NewLRU(size int) (*LRU, error) {
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &LRU{
		cache: c,
	}, nil
}

//Get get value by key
func (l *LRU) Get(key [32]byte) (interface{}, bool) {
	return l.cache.Get(key)
}

//Set set value to cache
func (l *LRU) Set(key [32]byte, val interface{}) {
	l.cache.Add(key, val)
}

//Delete delete value from cache
func (l *LRU) Delete(key [32]byte) {
	l.cache.Remove(key)
}
