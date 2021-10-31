package main

import (
	"flag"
	"fmt"
	"golearning/go-cache/lib"
	"time"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "c", "./config/n1_master.yaml", "your config file path")
	flag.Parse()
	err := lib.Boot(filePath)
	if err != nil {
		panic(err)
	}
	for true {
		time.Sleep(2 * time.Second)
		fmt.Println("当前leader：", lib.RaftNode.Leader())
	}
}
