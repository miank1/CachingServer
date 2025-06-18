package cache

import (
	"sync"
	"time"
)

type item struct {
	value      string
	expiration int64 // Unix timestamp in nanoseconds
}

type Cache struct {
	store  map[string]item
	mu     sync.RWMutex
	hits   int
	misses int
}

type Stats struct {
	Items  int `json:"items"`
	Hits   int `json:"hits"`
	Misses int `json:"misses"`
}

func NewCache() *Cache {
	c := &Cache{
		store: make(map[string]item),
	}

	go c.cleanupExpiredItems()

	return c
}

func (c *Cache) Set(key, value string, ttlSeconds int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(time.Duration(ttlSeconds) * time.Second).UnixNano()
	c.store[key] = item{value: value, expiration: expiration}

}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, found := c.store[key]
	if !found || time.Now().UnixNano() > it.expiration {
		c.misses++
		return "", false
	}

	if time.Now().UnixNano() > it.expiration {
		c.misses++
		return "", false
	}

	c.hits++
	return it.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
}

func (c *Cache) cleanupExpiredItems() {

	for {
		time.Sleep(1 * time.Second)
		c.mu.Lock()

		now := time.Now().UnixNano()

		for k, v := range c.store {
			if now > v.expiration {
				delete(c.store, k)
			}
		}

		c.mu.Unlock()

	}
}

func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Stats{
		Items:  len(c.store),
		Hits:   c.hits,
		Misses: c.misses,
	}
}
