package Janney

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	单元测试
*/

func TestGetterFunc_Get(t *testing.T) {
	g := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	b, _ := g.Get("name")
	assert.Equal(t, b, []byte("name"))
}

// 测试group的get方法
var db = map[string]string{
	"a": "aa",
	"b": "bb",
	"c": "cc",
}

func TestGroup_Get(t *testing.T) {
	// 记录key调用本地getter方法的次数
	cnt := make(map[string]int)

	g := NewGroup("test", 9, GetterFunc(func(key string) ([]byte, error) {
		cnt[key]++
		return []byte(db[key]), nil
	}))

	// 第一遍获取数据，缓存不在内存中，能从本地获取到数据
	// 第二遍获取数据，直接从缓存中获取
	for i := 0; i < 2; i ++ {
		for k, v := range db {
			bv, _ := g.Get(k)
			assert.Equal(t, v, bv.String())
			assert.Equal(t, 1, cnt[k])
		}
	}
}
