package ipc

import (
	"encoding/json"
	"jumie/internal/indexer"
	"log"
	"net"
)

type Server struct {
	socketPath string
	cache      *indexer.InfoCache
}

func NewServer(socketPath string, cache *indexer.InfoCache) *Server {
	return &Server{
		socketPath: socketPath,
		cache:      cache,
	}
}

func (s *Server) Listen() {
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("listener close error: %v\n", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer func(c net.Conn) {
				err := c.Close()
				if err != nil {
					log.Printf("close error: %v", err)
				}
			}(c)

			var data map[string]string
			err := json.NewDecoder(c).Decode(&data)
			if err != nil {
				log.Printf("error read data: %v", err)
			}

			log.Printf("received data %v\n", data)
		}(conn)
	}
}
