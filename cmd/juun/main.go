package main

import (
	"fmt"
	"juun/internal/indexer"
)

var globalCache = indexer.NewCache(indexer.SystemInfo{})

func main() {
	fmt.Println("juun client starting...")
	go indexer.RunIndexer(globalCache)
}
