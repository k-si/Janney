package Janney

import (
	"Janney/consistent_hash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
)

const (
	DefaultBasePath = "/janney/"
	replicas        = 50
)

type HttpPool struct {
	self       string // 自身服务端的地址和端口
	basePath   string // url增加一个字段标注该服务，比如127.0.0.1:8888/janney/xxx/xxx
	mu         sync.Mutex
	peers      *consistent_hash.Map  // 存储其他节点
	httpGetter map[string]PeerGetter // 记录节点地址和该节点getter映射
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:       self,
		basePath:   DefaultBasePath,
		peers:      consistent_hash.NewMap(replicas, nil),
		httpGetter: make(map[string]PeerGetter),
	}
}

// 提供本地节点的http服务，只接受/janney/group/key的请求
func (hp *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[Jenney receive:]", r.URL.Path)

	// 判断请求的base地址
	if !strings.HasPrefix(r.URL.Path, hp.basePath) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// 从url中抽取出group和key
	gk := strings.SplitN(r.URL.Path[len(hp.basePath):], "/", 2)
	if len(gk) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := gk[0]
	key := gk[1]

	// group get
	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, "no such group", http.StatusBadRequest)
		return
	}
	if len(key) == 0 {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
	view, err := g.Get(key)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 写回缓存数据
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(view.b); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

// peer是节点的地址，例如 127.0.0.1:8888/janney/
func (hp *HttpPool) Set(peers ...string) {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	// 将每个节点地址存入哈希环，并记录每个节点的getter
	for _, p := range peers {
		hp.peers.Add(p)
		hp.httpGetter[p] = &HttpGetter{baseURL: p + hp.basePath}
	}
}

// 获取key要请求的节点
func (hp *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	// 保证获取的地址不是空，而且不是自己
	if peer := hp.peers.Get(key); peer != "" && peer != hp.self {
		return hp.httpGetter[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HttpPool)(nil)

// 向其他节点发请求缓存数据的ds
type HttpGetter struct {
	baseURL string // 例如 127.0.0.1:8888/janney/
}

// 本地节点向别的节点请求缓存数据
func (hg *HttpGetter) Get(group, key string) ([]byte, error) {
	log.Println("[Janney] pick peer:", hg.baseURL)

	// 请求别的节点
	url := fmt.Sprintf("%v%v/%v", hg.baseURL, url2.QueryEscape(group), url2.QueryEscape(key))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 判断请求成功
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code:%v", resp.StatusCode)
	}

	// 读取body数据并返回
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// 确保*HttpGetter类型实现了PeerGetter接口
var _ PeerGetter = (*HttpGetter)(nil)
