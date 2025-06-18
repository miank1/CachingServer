package cache

type Cache struct {
	store map[string]string
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]string),
	}
}

func (c *Cache) Set(key, value string) {
	c.store[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	val, exist := c.store[key]
	return val, exist
}

func (c *Cache) Delete(key string) {
	delete(c.store, key)
}
