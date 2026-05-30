package main

import (
	"fmt"
	"jumie/internal/indexer"
)

var globalCache = indexer.NewCache(indexer.SystemInfo{})

func main() {
	go indexer.RunIndexer(globalCache)
	fmt.Println(globalCache)
}
