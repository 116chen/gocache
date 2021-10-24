package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (this *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := this.cache[key]; ok {
		this.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (this *Cache) RemoveOldest() {
	ele := this.ll.Back()
	if ele != nil {
		this.ll.Remove(ele)
		kv := ele.Value.(*entry)
		this.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		delete(this.cache, kv.key)
		if this.OnEvicted != nil {
			this.OnEvicted(kv.key, kv.value)
		}
	}
}

func (this *Cache) Add(key string, value Value) {
	if ele, ok := this.cache[key]; ok {
		this.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		this.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := this.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		this.cache[key] = ele
		this.nBytes += int64(value.Len()) + int64(len(key))
	}

	for this.maxBytes > 0 && this.nBytes > this.maxBytes {
		this.RemoveOldest()
	}
}

func (this *Cache) Len() int {
	return this.ll.Len()
}
