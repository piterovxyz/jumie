package main

import (
	"fmt"
	"jumie/internal/daemon"
	"jumie/internal/indexer"
	"jumie/internal/ipc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("jumie daemon starting...")

	_ = daemon.StartOllama()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nshutting down daemon...")
		daemon.StopOllama()
		os.Exit(0)
	}()

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
