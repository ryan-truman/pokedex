package pokecache

import (
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type cache struct {
	cacheEntries map[string]cacheEntry
}

func NewCache(time.Duration) {
}

func (c cache) Add(key string) {
	
}


func (c cache) Get(key string) []byte {
	return nil
}
