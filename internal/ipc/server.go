package ipc

import (
	"encoding/json"
	"jumie/internal/indexer"
	"log"
	"net"
	"os"
)

type Server struct {
	socketPath string
	cache      indexer.InfoCache
}

func (s *Server) Listen() {
	var path string

	if s.cache.Get().IsSU {
		path = "/var/jumie.sock"
	} else {
		path = os.Getenv("HOME") + "/.local/share/jumie/jumie.sock"
	}

	listener, err := net.Listen("unix", path)
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
