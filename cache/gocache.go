package cache

import (
	"fmt"
	pb "golearning/go-cache/cachepb"
	"golearning/go-cache/singalflight"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (this GetterFunc) Get(key string) ([]byte, error) {
	return this(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singalflight.Group
}

var (
	mu     sync.RWMutex
	once   sync.Once
	groups map[string]*Group
)

func NewGroup(name string, getter Getter, cacheBytes int64) *Group {
	if getter == nil {
		panic("getter nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if groups == nil {
		groups = make(map[string]*Group)
	}
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singalflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (this *Group) Get(key string) (ByteView, error) {
	// key is empty, return...
	if key == "" {
		return ByteView{}, fmt.Errorf("key is empty")
	}
	// get value from cache firstly...
	if value, ok := this.mainCache.get(key); ok {
		log.Println("cache hit. key =", key)
		return value, nil
	}
	// get value from db or remotes if cache is nonexistent...
	return this.load(key)
}

func (this *Group) load(key string) (ByteView, error) {
	view, err := this.loader.Do(key, func() (interface{}, error) {
		if this.peers != nil {
			if peer, ok := this.peers.PickPeer(key); ok {
				// 远端找不到，应该继续在本机找，而不是报错返回
				if value, err := this.getFromPeer(peer, key); err == nil {
					return value, nil
				}
			}
			log.Println("peer not found. key =", key)
		}
		return this.loadLocally(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return view.(ByteView), nil
}

func (this *Group) loadLocally(key string) (ByteView, error) {
	// use callback function to get value...
	tmpValue, err := this.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	view := ByteView{b: cloneByte(tmpValue)}
	// populate cache
	this.populateCache(key, view)
	return view, nil

}

func (this *Group) populateCache(key string, view ByteView) {
	this.mainCache.add(key, view)
}

func (this *Group) RegisterPeers(peers PeerPicker) {
	once.Do(func() {
		if peers == nil {
			panic("RegisterPeers only once, so peers can't be nil")
		}
		this.peers = peers
	})
}

func (this *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: this.name,
		Key:   key,
	}
	resp := &pb.Response{}
	err := peer.Get(req, resp)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: resp.Value}, nil
}
