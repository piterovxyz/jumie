package main

import (
	"fmt"
	"juun/internal/indexer"
	"log"
)

func main() {
	fmt.Println("juun client starting...")
	err := indexer.RunIndexer()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
