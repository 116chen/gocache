package lru

import "container/list"

// Cache LRU的实现逻辑
type Cache struct {
	// 容量
	maxBytes int64
	// 长度
	nBytes int64
	// 用于实现LRU的双向链表
	ll *list.List
	// 用于实现LRU的哈希表
	cache map[string]*list.Element
	// 回调函数，触发时机：移去最不常使用的节点的时候
	OnEvicted func(key string, value Value)
	// TODO 过期时间
}

// entry 缓存实例
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

// Get 获取缓存值
func (this *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := this.cache[key]; ok {
		// 将key对应enrty移到队头（表示最近使用）
		this.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移去最不常使用的节点
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

// Add 添加节点
func (this *Cache) Add(key string, value Value) {
	// 存在就移到队头
	if ele, ok := this.cache[key]; ok {
		this.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		this.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 不存在添加到队头
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
