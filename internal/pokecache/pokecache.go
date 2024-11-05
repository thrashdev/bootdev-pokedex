package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	values map[string]cacheEntry
	mu     *sync.RWMutex
}

func (cache *Cache) Add(key string, val []byte) {
	entry := cacheEntry{createdAt: time.Now(), val: val}
	cache.values[key] = entry
}

func (cache *Cache) Get(key string) (val []byte, success bool) {
	entry, success := cache.values[key]
	if success {
		return entry.val, success
	}
	return []byte{}, false
}

func (cache *Cache) reapLoop(ticker *time.Ticker, interval time.Duration) {
	for tick := range ticker.C {
		cache.mu.Lock()
		for k, v := range cache.values {
			if tick.Sub(v.createdAt) > interval {
				delete(cache.values, k)
			}
		}
		cache.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	ticker := time.NewTicker(interval)
	cache := Cache{values: make(map[string]cacheEntry), mu: &sync.RWMutex{}}
	go cache.reapLoop(ticker, interval)
	return &cache

}
