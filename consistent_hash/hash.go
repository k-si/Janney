package hash

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

func NewMap(hash HashFunc, replicas int) *Map {
	m := &Map{
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

func (m *Map) Remove(nodes ...string) {
	for _, n := range nodes {
		for i := 0; i < m.replicas; i ++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + n)))
			delete(m.vnode, h)
			idx := sort.SearchInts(m.sortedKeys, h)
			m.sortedKeys = append(m.sortedKeys[:idx], m.sortedKeys[idx+1:]...)
		}
	}
	sort.Ints(m.sortedKeys)
}

func (m *Map) Get(node string) string {
	
}
