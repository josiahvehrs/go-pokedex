package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	Entries  map[string]cacheEntry
	Interval time.Duration
	mu       sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{Interval: interval, Entries: make(map[string]cacheEntry), mu: sync.RWMutex{}}
	cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{createdAt: time.Now(), val: val}
	c.Entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.Entries[key]
	if !ok {
		return nil, ok
	}

	return entry.val, ok
}

func (c *Cache) reapLoop() {
	tickChan := time.NewTicker(c.Interval)
	go func() {
		for range tickChan.C {
			c.mu.Lock()
			for key, val := range c.Entries {
				if time.Since(val.createdAt) > c.Interval {
					delete(c.Entries, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}
