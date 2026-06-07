package ipc

import (
	"context"
	"encoding/json"
	"fmt"
	"jumie/internal/ai"
	"jumie/internal/config"
	"jumie/internal/indexer"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
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

	dir := filepath.Dir(s.socketPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create socket directory: %v\n", err)
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

		go func(c net.Conn) {
			req := new(Request)

			err := json.NewDecoder(c).Decode(&req)
			if err != nil {
				fmt.Printf("data decode error: %v", err)
				return
			}

			switch req.Type {
			case "plan":
				s.doPlan(req.AIMessage, c)
			case "exec":
				s.doExec(req.Commands, c)
			}
		}(conn)
	}
}

func (s *Server) doPlan(msg string, c net.Conn) {
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			log.Printf("close error: %v", err)
		}
	}(c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error load config: %v\n", err)
	}

	client, err := ai.NewClient(cfg)
	if err != nil {
		fmt.Printf("ai error: %v\n", err)
		return
	}

	err = client.UpdateCache(s.cache.Get())
	if err != nil {
		fmt.Printf("update cache error: %v", err)
		return
	}

	plan, err := client.GeneratePlan(ctx, msg)

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
}

func (s *Server) doExec(commands []string, c net.Conn) {
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			log.Printf("close error: %v", err)
		}
	}(c)

	shell := s.cache.Get().Shell
	if shell == "" {
		shell = "/bin/sh"
	}

	for _, cmd := range commands {
		_, err := fmt.Fprintf(c, "\nrunning: %s\n", cmd)
		if err != nil {
			return
		}

		cmd := exec.Command(shell, "-c", cmd)
		cmd.Stdout = c
		cmd.Stderr = c

		if err := cmd.Run(); err != nil {
			_, err := fmt.Fprintf(c, "error executing command: %v\n", err)
			if err != nil {
				return
			}
			return
		}
	}

	_, err := fmt.Fprintln(c, "\nsuccessfully executed all commands!")
	if err != nil {
		return
	}
}
