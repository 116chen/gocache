package consistenthash

import (
	"strconv"
	"testing"
)

func TestNewMap(t *testing.T) {
	hash := NewMap(func(key []byte) uint32 {
		val, _ := strconv.Atoi(string(key))
		return uint32(val)
	}, 3)
	// 2 4 6 12 14 16 22 24 26
	hash.Add("2", "6", "4")
	cases := map[string]string{
		"11": "2",
		"13": "4",
		"17": "2",
		"27": "2",
	}
	for k, v := range cases {
		if hash.Get(k) != v {
			t.Fatal("case1 not pass")
		}
	}

	// 2 4 6 8 12 14 16 18 22 24 26 28
	hash.Add("8")
	cases["17"] = "8"
	cases["27"] = "8"
	for k, v := range cases {
		if hash.Get(k) != v {
			t.Fatal("case2 not pass")
		}
	}
}
