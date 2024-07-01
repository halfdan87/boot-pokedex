package pokecache

import (
    "fmt"
    "sync"
    "time"
)

type Cache struct {
    mu sync.RWMutex
    items map[string]cacheEntry
}

func NewCache(dur time.Duration) *Cache {
    cache := &Cache{
        mu: sync.RWMutex{},
        items : make(map[string]cacheEntry),
    }

    go cache.reapLoop(dur)

    return cache
 }

type cacheEntry struct {
    createdAt time.Time
    val []byte
}

func (c *Cache) Add(key string, newVal []byte) {
    c.mu.Lock()
    c.items[key] = cacheEntry{
        createdAt : time.Now(),
        val : newVal,
    }
    c.mu.Unlock()
}

func (c *Cache) reapLoop(dur time.Duration) {
    ticker := time.NewTicker(dur)

    defer ticker.Stop()

    for range ticker.C {
        fmt.Println("I work")

        keysToDelete := []string{}
        for key, val := range c.items {
            if time.Since(val.createdAt) > dur {
                keysToDelete = append(keysToDelete, key)
            }
        }
        for _, key := range keysToDelete {
            c.mu.Lock()
            delete(c.items, key)
            c.mu.Unlock()
        }
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    if entry, ok := c.items[key]; ok {
        return entry.val, true
    }
    return nil, false
}

