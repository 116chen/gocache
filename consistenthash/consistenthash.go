package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(key []byte) uint32

type Map struct {
	// 自定义hash函数
	hash Hash
	// 每个实例节点对应虚拟节点个数（包括自己）
	replicas int
	// 所有节点的hashCode
	keys []int
	// 虚拟节点{hashCode}---->实例节点的{ip:port}
	hashMap map[int]string
}

func NewMap(hash Hash, replicas int) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加虚拟节点
func (this *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < this.replicas; i++ {
			hashCode := int(this.hash([]byte(strconv.Itoa(i) + key)))
			this.keys = append(this.keys, hashCode)
			this.hashMap[hashCode] = key
		}
	}
	sort.Ints(this.keys)
}

// Get 采用哈希一致性获取实例节点的{ip:port}
func (this *Map) Get(key string) string {
	if len(this.keys) == 0 {
		return ""
	}
	hashCode := int(this.hash([]byte(key)))
	idx := sort.Search(len(this.keys), func(i int) bool {
		return this.keys[i] > hashCode
	})
	// cycle... if idx == len(this.keys) then 0
	return this.hashMap[this.keys[idx%len(this.keys)]]
}
