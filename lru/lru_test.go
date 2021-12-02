package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type String string
func (s String) Len() int64 {
	return int64(len(s))
}

/*
	单元测试
 */
func TestCache_Get(t *testing.T) {
	c := NewCache(1024, nil)
	v, ok := c.Get("name")
	assert.Nil(t, v)
	assert.False(t, ok)
}

func TestCache_Put(t *testing.T) {
	c := NewCache(1024, nil)
	c.Put("name", String("zhang san"))
}

func TestCache_RemoveOldest(t *testing.T) {
	c := NewCache(1024, nil)
	c.RemoveOldest()
}

func TestCache_OnEvicted(t *testing.T) {
	flag := false
	onEvicted := func(key string, value Value) {
		flag = true
	}
	c := NewCache(1024, onEvicted)
	c.Put("name", String("zhang san"))
	c.RemoveOldest()
	assert.True(t, flag)
}

/*
	场景测试
 */

// 测试组合get put
func TestCache_1(t *testing.T) {
	c := NewCache(1024, nil)
	c.Put("name", String("zhang san"))
	name, ok := c.Get("name")
	assert.Equal(t, String("zhang san"), name)
	assert.True(t, ok)

	c.Put("name", String("li si"))
	name, ok = c.Get("name")
	assert.Equal(t, String("li si"), name)
	assert.True(t, ok)
}

// 测试缓存淘汰
func TestCache_2(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"

	size := len(k1) + len(k2) + len(v1) + len(v2)
	c := NewCache(int64(size), nil)
	c.Put(k1, String(v1))
	c.Put(k2, String(v2))
	c.Put(k3, String(v3))
	v, ok := c.Get(k1)
	assert.Nil(t, v)
	assert.False(t, ok)
}
