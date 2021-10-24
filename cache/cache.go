package cache

import (
	"golearning/go-cache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (this *cache) add(key string, value ByteView) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.lru == nil {
		this.lru = lru.NewCache(this.cacheBytes, nil)
	}
	this.lru.Add(key, value)
}

func (this *cache) get(key string) (value ByteView, ok bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.lru == nil {
		return
	}
	if ele, ok := this.lru.Get(key); ok {
		return ele.(ByteView), ok
	}
	return
}
