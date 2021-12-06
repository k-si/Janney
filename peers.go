package Janney

// PeerGetter可以通过group和key来获取缓存值
// 其实就是通过http请求获取对应节点的缓存值
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}

// 可以通过key来获取PeerGetter
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}


