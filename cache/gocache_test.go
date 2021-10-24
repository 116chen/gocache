package cache

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGroup_Get(t *testing.T) {
	db := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
	}
	count := make(map[string]int)
	f := GetterFunc(func(key string) ([]byte, error) {
		log.Println("slowDb load. key=", key)
		if _, ok := count[key]; !ok {
			count[key] = 0
		}
		count[key]++
		if value, ok := db[key]; ok {
			return []byte(value), nil
		}
		return nil, fmt.Errorf("key is not found")
	})

	group := NewGroup("ch", f, 1<<20)

	for k, v := range db {
		if value, err := group.Get(k); err != nil || count[k] != 1 || value.String() != v {
			t.Fatal("case1 failed")
		}
		if value, err := group.Get(k); err != nil || count[k] != 1 || value.String() != v {
			t.Fatal("case2 failed")
		}
	}
	t.Log("case pass")
}

func TestOutFile(t *testing.T) {
	file, err := os.Open("./logs.txt")
	if err != nil {
		panic("open file err")
	}
	ans, err := os.Create("./ans.csv")
	defer func() {
		file.Close()
		ans.Close()
	}()
	buf := bufio.NewReader(file)
	f := csv.NewWriter(ans)
	res := make([][]string, 0)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		iStart := strings.Index(line, "InterceptedNo")
		iEnd := strings.Index(line, "req.GetShopId")
		interceptedNo := line[iStart+len("InterceptedNo")+1 : iEnd-1]

		iStart = strings.Index(line, "GetShopId")
		iEnd = strings.Index(line, "_msg=")
		shopId := line[iStart+len("GetShopId")+1 : iEnd-1]
		fmt.Println(interceptedNo, shopId)
		res = append(res, []string{shopId, interceptedNo})
		if err != nil {
			if err == io.EOF {
				fmt.Println("file read ok")
				break
			} else {
				return
			}
		}
	}
	f.WriteAll(res)
}
