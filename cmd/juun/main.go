package main

import (
	"fmt"
	"juun/internal/indexer"
	"log"
)

func main() {
	var sysInfo indexer.InfoCache

	fmt.Println("juun client starting...")
	sysInfo, err := indexer.RunIndexer()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
