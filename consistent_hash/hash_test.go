package consistent_hash

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestMap_Add(t *testing.T) {
	m := NewMap(3, func(bytes []byte) uint32 {
		h, _ := strconv.Atoi(string(bytes))
		return uint32(h)
	})
	m.Add("1")
}

func TestMap_Get(t *testing.T) {
	m := NewMap(3, func(bytes []byte) uint32 {
		h, _ := strconv.Atoi(string(bytes))
		return uint32(h)
	})
	node := m.Get([]byte("1"))
	assert.Equal(t, node, "")
}

func TestMap_Remove(t *testing.T) {
	m := NewMap(3, func(bytes []byte) uint32 {
		h, _ := strconv.Atoi(string(bytes))
		return uint32(h)
	})
	m.Remove("1")
}

func TestMap(t *testing.T) {
	// 这里保证存入的值都是数字，哈希结果就是数字本身的值
	m := NewMap(3, func(bytes []byte) uint32 {
		h, _ := strconv.Atoi(string(bytes))
		return uint32(h)
	})

	// 节点哈希值为：02 12 22 04 14 24 06 16 26
	// sortedKeys: 02 04 06 12 14 16 22 24 26
	m.Add("2", "4", "6")

	assert.Equal(t, "2", m.Get([]byte("0")))
	assert.Equal(t, "4", m.Get([]byte("3")))
	assert.Equal(t, "6", m.Get([]byte("25")))
	assert.Equal(t, "2", m.Get([]byte("27")))

	// sortedKeys: 04 06 14 16 24 26
	m.Remove("2")
	assert.Equal(t, "4", m.Get([]byte("0")))
	assert.Equal(t, "4", m.Get([]byte("3")))
	assert.Equal(t, "6", m.Get([]byte("25")))
	assert.Equal(t, "4", m.Get([]byte("27")))

	m.Remove("4", "6")
	assert.Equal(t, "", m.Get([]byte("0")))
}