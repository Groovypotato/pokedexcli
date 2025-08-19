package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mutex sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache {
	var newCache Cache
	newCache.entries = make(map[string]cacheEntry)
	newCache.interval = interval
	go (&newCache).reapLoop()
	return  &newCache
}

func (c *Cache) Add( key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	newEntry := cacheEntry {
		createdAt: time.Now(),
		val: val,
	}
	c.entries[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop () {
	ticker := time.Tick(c.interval)
	for {
		<-ticker
		c.mutex.Lock()
		oldTime := time.Now().Add(-c.interval)
		for key, entry := range c.entries {
			if entry.createdAt.Before(oldTime) {
				delete(c.entries,key)
			}
		}
		c.mutex.Unlock()
	}

}