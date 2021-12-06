package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func([]byte) uint32

type Map struct {
	hash       HashFunc       // 求服务节点哈希值的算法
	sortedKeys []int          // 哈希环，存储服务节点的哈希值
	vnode      map[int]string // 虚拟节点映射
	replicas   int            // 虚拟节点扩容倍数
}

func NewMap(replicas int, hash HashFunc) *Map {
	m := &Map{
		hash:       hash,
		sortedKeys: make([]int, 0),
		vnode:      make(map[int]string),
		replicas:   replicas,
	}
	if hash == nil {
		m.hash = crc32.ChecksumIEEE // 默认使用crc
	}
	return m
}

// 添加n个服务节点
func (m *Map) Add(nodes ...string) {
	for _, n := range nodes {
		for i := 0; i < m.replicas; i++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + n))) // hashcode -> 0servername
			m.sortedKeys = append(m.sortedKeys, h)
			m.vnode[h] = n
		}
	}
	sort.Ints(m.sortedKeys)
}

// 删除服务节点
func (m *Map) Remove(nodes ...string) {
	for _, n := range nodes {
		for i := 0; i < m.replicas; i++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + n)))
			idx := sort.SearchInts(m.sortedKeys, h)

			// 如果在哈希环找不到要删除的节点，可直接结束寻找
			if idx == len(m.sortedKeys) || m.sortedKeys[idx] != h {
				break
			}
			m.sortedKeys = append(m.sortedKeys[:idx], m.sortedKeys[idx+1:]...)
			delete(m.vnode, h)
		}
	}
	sort.Ints(m.sortedKeys)
}

// 返回key应当访问的服务节点
func (m *Map) Get(key string) string {
	if len(key) == 0 || len(m.sortedKeys) == 0 {
		return ""
	}
	h := int(m.hash([]byte(key)))

	// 当找不到h时，返回数组的长度
	idx := sort.Search(len(m.sortedKeys), func(i int) bool {
		return m.sortedKeys[i] >= h
	})
	return m.vnode[m.sortedKeys[idx%len(m.sortedKeys)]]
}
