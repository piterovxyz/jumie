package main

import (
	"encoding/json"
	"jumie/internal/ipc"
	"log"
	"net"
	"os"
	"os/user"
)

func main() {
	var path string

	current, err := user.Current()
	if err != nil {
		log.Fatalf("error to get current user: %v\n", err)
	}

	if current.Uid == "0" {
		path = "/var/jumie.sock"
	} else {
		path = os.Getenv("HOME") + "/.local/share/jumie/jumie.sock"
	}

	c, err := ipc.NewClient(path)
	if err != nil {
		log.Fatalf("error creating ipc client: %v\n", err)
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v\n", err)
		}
	}(c.Conn)

	data := ipc.Request{
		Payload: json.RawMessage(`{"ai_message": "install zsh and oh-my-zsh"}`),
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("failed to marshal json: %v", err)
	}

	_, err = c.Conn.Write(bytes)
	if err != nil {
		log.Fatalf("failed to write to socket: %v", err)
	}

	var resp ipc.Response
	err = json.NewDecoder(c.Conn).Decode(&resp)
	if err != nil {
		log.Fatalf("read error: %v", err)
	}

	log.Printf("received response from daemon: %v", resp)
}
