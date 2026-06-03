package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/user"
	"time"
)

func main() {
	for {
		var path string

		current, err := user.Current()
		if err != nil {
			log.Printf("error to get current user: %v\n", err)
			return
		}

		if current.Uid == "0" {
			path = "/var/jumie.sock"
		} else {
			path = os.Getenv("HOME") + "/.local/share/jumie/jumie.sock"
		}

		c, err := net.Dial("unix", path)
		if err != nil {
			log.Printf("failed to connect socket: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}

		data := map[string]string{
			"name": "artur",
			"age":  "21",
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("failed to marshal json: %v", err)
		}

		_, err = c.Write(bytes)
		if err != nil {
			log.Printf("failed to write to socket: %v", err)
		}

		time.Sleep(time.Second * 5)
	}
}
