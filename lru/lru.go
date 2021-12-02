package lru

import (
	"container/list"
)

/*
	实现LRU缓存淘汰机制
*/
type Cache struct {
	maxBytes int64 // 允许的最大内存
	nBytes   int64 // 当前占用的内存
	ll       *list.List
	cache    map[string]*list.Element // element存储*entry类型

	// 缓存淘汰回调函数
	OnEvicted func(key string, value Value)
}

// 将数据的key作为附加值一起缓存起来
type entry struct {
	key   string
	value Value
}

// 实际要缓存的对象必须实现Value接口
type Value interface {
	Len() int64 // 获取单个缓存数据所占字节个数
}

func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) GetLL() *list.List {
	return c.ll
}

// 命中缓存
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.ll.MoveToFront(ele) // 调整缓存节点到队尾
		return kv.value, true
	}
	return nil, false
}

// 覆盖写缓存
func (c *Cache) Put(key string, value Value) {

	// 如果加入后超出缓存最大限制，需要先淘汰一部分缓存
	size := int64(len(key)) + value.Len()
	for c.maxBytes != 0 && c.nBytes+size > c.maxBytes {
		c.RemoveOldest()
	}

	// kv覆盖写入
	if ele, ok := c.cache[key]; ok {

		// 已经存在相同key，将缓存节点的value更新
		kv := ele.Value.(*entry)
		kv.value = value

		// 缓存节点移至队尾
		c.ll.MoveToFront(ele)

		// 更新已使用的内存容量
		c.nBytes += value.Len() - kv.value.Len()

	} else {

		// 新的kv，将缓存节点加入队尾
		ele := c.ll.PushFront(&entry{key: key, value: value})

		// 更新索引
		c.cache[key] = ele

		// 更新已使用的内存容量
		c.nBytes += size
	}
}

// 淘汰缓存
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {

		// 删除缓存节点
		c.ll.Remove(ele)

		// 删除缓存索引
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)

		// 更新已使用的内存容量
		c.nBytes -= int64(len(kv.key)) + kv.value.Len()

		// 回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 缓存节点个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
