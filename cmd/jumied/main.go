package main

import (
	"fmt"
	"jumie/internal/indexer"
	"jumie/internal/ipc"
)

func main() {
	info := indexer.SystemInfo{}
	cache := indexer.NewCache(info)
	indexer.RunIndexer(cache)
	ipc.Listen()

	fmt.Println("jumie daemon starting...")
}
