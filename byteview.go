package Janney

// 实际作为缓存的数据结构，只支持读
type ByteView struct {
	b []byte
}

// 实际缓存的数据要实现Value接口
func (bv ByteView) Len() int64 {
	return int64(len(bv.b))
}

// 返回缓存字节数组的副本
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

// 将缓存的字节数组转为string
func (bv ByteView) String() string {
	return string(bv.b)
}

func cloneBytes(b []byte) []byte {
	bb := make([]byte, len(b))
	copy(bb, b)
	return bb
}