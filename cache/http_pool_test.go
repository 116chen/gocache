package cache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestNewHttpPool(t *testing.T) {
	db := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
	}
	NewGroup("test", GetterFunc(func(key string) ([]byte, error) {
		log.Println("slowDb load. key=", key)
		if value, ok := db[key]; ok {
			return []byte(value), nil
		}
		return nil, fmt.Errorf("key is not found")
	}), 1<<20)

	host := "localhost:9999"
	peer := NewHttpPool(host)
	_ = http.ListenAndServe(host, peer)
}
