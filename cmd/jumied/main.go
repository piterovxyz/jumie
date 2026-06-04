package main

import (
	"fmt"
	"jumie/internal/indexer"
	"jumie/internal/ipc"
	"os"
)

func main() {
	fmt.Println("jumie daemon starting...")

	var server *ipc.Server
	info := indexer.SystemInfo{}
	cache := indexer.NewCache(info)
	indexer.RunIndexer(cache)

	if cache.Get().IsSU {
		server = ipc.NewServer("/var/jumie.sock", cache)
	} else {
		server = ipc.NewServer(os.Getenv("HOME")+"/.local/share/jumie/jumie.sock", cache)
	}

	server.Listen()
}
