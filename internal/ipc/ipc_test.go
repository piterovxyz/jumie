package ipc

import (
	"encoding/json"
	"jumie/internal/ai"
	"jumie/internal/indexer"
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(tempDir, ".local", "share", "jumie", "jumie.sock")
	if client.socketPath != expected {
		t.Errorf("expected %q, got %q", expected, client.socketPath)
	}
}

func TestClientPing(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "test.sock")

	l, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		var req Request
		json.NewDecoder(conn).Decode(&req)
		if req.Type == "ping" {
			conn.Write([]byte(`{"status":"ok"}`))
		}
	}()

	client := &Client{socketPath: sockPath}
	err = client.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientDoPlan(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "test2.sock")

	l, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		var req Request
		json.NewDecoder(conn).Decode(&req)
		if req.Type == "exec" && len(req.Commands) == 1 && req.Commands[0] == "echo 1" {
			conn.Write([]byte("ok"))
		}
	}()

	client := &Client{socketPath: sockPath}
	plan := &ai.Plan{
		Steps: []ai.Step{{Command: "echo 1"}},
	}
	err = client.DoPlan(plan)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerPingAndExec(t *testing.T) {
	tempDir := t.TempDir()
	sockPath := filepath.Join(tempDir, "server.sock")

	cache := indexer.NewCache(indexer.SystemInfo{Shell: "/bin/sh"})

	server := NewServer(sockPath, cache)
	go server.Listen()

	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("unix", sockPath)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	req := Request{Type: "ping"}
	json.NewEncoder(conn).Encode(req)

	var res map[string]string
	err = json.NewDecoder(conn).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}
	if res["status"] != "ok" {
		t.Errorf("expected status ok, got %v", res["status"])
	}

	conn2, err := net.Dial("unix", sockPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()

	req2 := Request{Type: "exec", Commands: []string{"echo test"}}
	json.NewEncoder(conn2).Encode(req2)

	buf := make([]byte, 1024)
	n, _ := conn2.Read(buf)
	if n == 0 {
		t.Errorf("expected output from exec")
	}
}
