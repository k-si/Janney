package Janney

import (
	"errors"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// 当缓存不存在时的回调函数
type GetterFunc func(key string) ([]byte, error)

func (gf GetterFunc) Get(key string) ([]byte, error) {
	return gf(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name      string
	getter    Getter // 用户自定义方法，用于获取不在内存中的数据
	mainCache *Cache // 带有并发控制的cache，LRU淘汰策略
}

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("getter can not be nil")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: &Cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) (*Group) {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if "" == key {
		return ByteView{}, errors.New("key can not be empty string")
	}
	if value, ok := g.mainCache.Get(key); ok {
		log.Println("[Jenney hit]:", string(value.b))
		return value, nil
	}

	// 没有命中缓存，从别处调到缓存中
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// 从本地获取数据
func (g *Group) getLocally(key string) (ByteView, error) {

	// 用户自定义从本地获取数据的方式
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	// 保证ByteView不可变，写入副本值
	value := ByteView{b: cloneBytes(bytes)}

	// 加载到缓存
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Put(key, value)
}
