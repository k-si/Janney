package Janney

import (
	"log"
	"net/http"
	"strings"
)

const DefaultBasePath = "/janney/"

type HttpPool struct {
	self string // 服务端的地址和端口
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: DefaultBasePath,
	}
}

func (hp *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[Jenney receive:]", r.URL.Path)

	// 判断请求的base地址
	if !strings.HasPrefix(r.URL.Path, hp.basePath) {
		http.Error(w, "bad request", 400)
		return
	}

	// 从url中抽取出group和key
	gk := strings.SplitN(r.URL.Path[len(hp.basePath):], "/", 2)
	if len(gk) != 2 {
		http.Error(w, "bad request", 400)
		return
	}
	groupName := gk[0]
	key := gk[1]

	// group get
	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, "no such group", 400)
		return
	}
	if len(key) == 0 {
		http.Error(w, "missing key", 400)
		return
	}
	view, err := g.Get(key)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}

	// 写回缓存数据
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(view.ByteSlice()); err != nil {
		http.Error(w, "internal error", 500)
		return
	}
}
