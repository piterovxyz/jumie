package main

import (
	"fmt"
	"jumie/internal/indexer"
	"jumie/internal/ipc"
)

func main() {
	fmt.Println("jumie daemon starting...")

	info := indexer.SystemInfo{}
	cache := indexer.NewCache(info)
	indexer.RunIndexer(cache)
	ipc.Listen(&info)
}
