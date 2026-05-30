package main

import (
	"fmt"
	"jumie/internal/indexer"
)

func main() {
	info := indexer.SystemInfo{}
	cache := indexer.NewCache(info)
	indexer.RunIndexer(cache)

	fmt.Println("jumie daemon starting...")
}
