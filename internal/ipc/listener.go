package ipc

import (
	"encoding/json"
	"jumie/internal/indexer"
	"log"
	"net"
	"os"
)

func Listen(sysInfo indexer.SystemInfo) {
	var path string

	if sysInfo.IsSU {
		path = "/var/jumie.sock"
	} else {
		path = os.Getenv("HOME") + "/.local/share/jumie/jumie.sock"
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			var data map[string]string
			err := json.NewDecoder(c).Decode(&data)
			if err != nil {
				log.Printf("error read data: %v", err)
			}

			log.Printf("received data %v\n", data)
		}(conn)
	}
}
