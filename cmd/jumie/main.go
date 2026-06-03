package main

import (
	"jumie/internal/ipc"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <message>", os.Args[0])
	}

	msg := os.Args[1]

	c, err := ipc.NewClient()
	if err != nil {
		log.Fatalf("error creating ipc client: %v\n", err)
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("error closing connection: %v\n", err)
		}
	}(c.Conn)

	err = c.SendMessage(msg)
	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
	}
}
