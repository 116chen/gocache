package lru

import (
	"testing"
)

type String string

func (this String) Len() int {
	return len(this)
}

func TestCache_Get(t *testing.T) {
	cache := NewCache(0, nil)
	cache.Add("key1", String("1234"))
	if ele, ok := cache.Get("key1"); !ok || string(ele.(String)) != "1234" {
		t.Fatal("not pass")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Fatal("not pass")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cache := NewCache(int64(len(k1+k2+v1+v2)), nil)
	cache.Add(k1, String(v1))
	if ele, ok := cache.Get(k1); !ok || string(ele.(String)) != v1 {
		t.Fatal("case1 not pass")
	}
	cache.Add(k2, String(v2))
	cache.Add(k3, String(v3))
	if _, ok := cache.Get(k1); ok {
		t.Fatal("case2 not pass")
	}
}

func TestOnEvicted(t *testing.T) {

}
