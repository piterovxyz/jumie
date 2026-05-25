package main

import (
	"fmt"
	"juun/internal/indexer"
)

func main() {
	fmt.Println("juun client starting...")
	indexer.RunIndexer()
}
