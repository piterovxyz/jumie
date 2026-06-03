package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

func main() {
	for {
		c, err := net.Dial("unix", "/var/jumie.sock")
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
