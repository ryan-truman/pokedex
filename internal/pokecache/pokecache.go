package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cacheEntries map[string]cacheEntry
	interval time.Duration
	ticker *time.Ticker
	mu sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	pokeCache := Cache {
		interval: interval,
		ticker: time.NewTicker(interval), 
	}
	go pokeCache.reapLoop()
	return &pokeCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEntries[key]= cacheEntry {
		createdAt: time.Now(),
		val: val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cacheEntries[key]
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	for  range c.ticker.C {
		c.mu.Lock()
		for key, entry := range c.cacheEntries {
			if entry.createdAt.Before(time.Now().Add(-c.interval)) {
				delete(c.cacheEntries, key)
			}
		}
		c.mu.Unlock()
	}
}
