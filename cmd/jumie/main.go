package main

import (
	"fmt"
	"jumie/internal/indexer"
	"log"
)

var globalCache = indexer.NewCache(indexer.SystemInfo{})

func main() {
	var sysInfo *indexer.InfoCache

	fmt.Println("jumie client starting...")
	sysInfo, err := indexer.RunIndexer()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(sysInfo)
}
