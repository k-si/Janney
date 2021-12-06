package main

import (
	"Janney"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *Janney.Group {
	g := Janney.NewGroup("score", 1024, Janney.GetterFunc(func(key string) ([]byte, error) {
		v, ok := db[key]
		if ok {
			return []byte(v), nil
		}
		return []byte{}, errors.New("no that key")
	}))
	return g
}

func startCacheServe(addr string, peersAddr []string, group *Janney.Group) {
	hp := Janney.NewHttpPool(addr)
	hp.Set(peersAddr...)
	group.RegistryPeers(hp)
	log.Println("cache service start at:", addr)
	http.ListenAndServe(addr[7:], hp)
}

func startApiServe(addr string, group *Janney.Group) {
	http.Handle("/api", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		k := request.URL.Query().Get("key")
		bv, err := group.Get(k)
		if err != nil {
			http.Error(writer, "internal error", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Write(bv.ByteSlice())
	}))
	log.Println("api service start at:", addr)
	http.ListenAndServe(addr[7:], nil)
}

func main() {
	var port int
	var api bool
	var apiPort int
	flag.IntVar(&port, "port", 8001, "cache server port")
	flag.BoolVar(&api, "api", false, "open api service")
	flag.IntVar(&apiPort, "apiPort", 9999, "api service port")
	flag.Parse()

	jan := createGroup()

	peersMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var peers []string
	for _, v := range peersMap {
		peers = append(peers, v)
	}

	if api {
		go startApiServe(fmt.Sprintf("http://localhost:%v", apiPort), jan)
	}

	startCacheServe(peersMap[port], peers, jan)
}
