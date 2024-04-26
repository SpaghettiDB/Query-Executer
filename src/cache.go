package executer

import "sync"

type Cache struct {
	name string
	data map[string]interface{}
	mux sync.Mutex
}

func NewCache(name string) *Cache {
	return &Cache{
		name: name,
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) interface{} {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.data[key]
}

func (c *Cache) Contains(key string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.data[key]
	return ok
}

func (c *Cache) Invalidate(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.data, key)
}

func (c *Cache) InvalidateAll() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.data = make(map[string]interface{})
}

func (c *Cache) Size() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.data)
}


func (c *Cache) Sync (){
	c.mux.Lock()
	defer c.mux.Unlock()
	// sync the cache with the storage engine
	// for example, write all the cache data to the disk
}