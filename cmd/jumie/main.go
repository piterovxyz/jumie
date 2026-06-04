package main

import (
	"fmt"
	"jumie/internal/ipc"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <message>", os.Args[0])
	}

	msg := strings.Join(os.Args[1:], " ")

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

	stop := startSpinner()
	resp, err := c.SendMessage(msg)
	stop()

	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
	}

	fmt.Println(resp)
}
