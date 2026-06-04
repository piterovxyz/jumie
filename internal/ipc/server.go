package ipc

import (
	"context"
	"encoding/json"
	"fmt"
	"jumie/internal/ai"
	"jumie/internal/indexer"
	"log"
	"net"
	"os"
	"time"
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
	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		log.Printf("error removing socket: %v", err)
	}

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

		go doRequest(conn)
	}
}

func doRequest(c net.Conn) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			log.Printf("close error: %v", err)
		}
	}(c)

	var req Request
	err := json.NewDecoder(c).Decode(&req)
	if err != nil {
		fmt.Printf("data decode error: %v", err)
		return
	}

	msg, err := req.Payload.MarshalJSON()
	if err != nil {
		fmt.Printf("payload decode error: %v", err)
		return
	}
	log.Printf("received data %s\n", msg)

	client, err := ai.NewClient(os.Getenv("GEMINI_API_KEY"), os.Getenv("GEMINI_MODEL"))
	if err != nil {
		fmt.Printf("ai error: %v\n", err)
		return
	}

	plan, err := client.GeneratePlan(ctx, string(msg))
	if err != nil {
		fmt.Printf("plan error: %v\n", err)
		return
	}

	data, err := json.Marshal(plan)

	_, err = c.Write(data)
	if err != nil {
		log.Printf("write error: %v", err)
		return
	}

	<-ctx.Done()
}
