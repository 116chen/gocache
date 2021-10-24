package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(key []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
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
