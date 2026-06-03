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

		user, err := user.Current()
		if err != nil {
			log.Printf("error to get current user: %v\n", err)
			return
		}

		if user.Uid == "0" {
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

		json, err := json.Marshal(data)
		if err != nil {
			log.Printf("failed to marshal json: %v", err)
		}

		c.Write(json)
		time.Sleep(time.Second * 5)
	}
}
