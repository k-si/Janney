package Janney

import (
	"Janney/lru"
	"sync"
)

// 对缓存做并发控制
type Cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

func (c *Cache) Put(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 延迟创建lru
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Put(key, value)
}
