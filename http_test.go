package Janney

import (
	"errors"
	"log"
	"net/http"
	"testing"
)

func TestHttpPool_ServeHTTP(t *testing.T) {
	db := map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
	}
	NewGroup("test", 1024, GetterFunc(func(key string) ([]byte, error) {
		v, ok := db[key]
		if !ok {
			return []byte{}, errors.New("no such key")
		}
		log.Println("[from slow DB]:", v)
		return []byte(v), nil
	}))

	addr := "localhost:9999"
	hp := NewHttpPool(addr)
	_ = http.ListenAndServe(addr, hp)
}
