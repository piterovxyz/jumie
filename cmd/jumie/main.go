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
	for {
		err := func() error {
			var path string

			current, err := user.Current()
			if err != nil {
				log.Printf("error to get current user: %v\n", err)
				return err
			}

			if current.Uid == "0" {
				path = "/var/jumie.sock"
			} else {
				path = os.Getenv("HOME") + "/.local/share/jumie/jumie.sock"
			}

			c, err := ipc.NewClient(path)
			if err != nil {
				log.Printf("error creating ipc client: %v\n", err)
				return err
			}
			defer func(Conn net.Conn) {
				err := Conn.Close()
				if err != nil {
					log.Printf("error closing connection: %v\n", err)
				}
			}(c.Conn)

			data := ipc.Request{
				Payload: json.RawMessage(`{"ai_message": "install zsh and oh-my-zsh"}`),
			}

			bytes, err := json.Marshal(data)
			if err != nil {
				log.Printf("failed to marshal json: %v", err)
			}

			_, err = c.Conn.Write(bytes)
			if err != nil {
				log.Printf("failed to write to socket: %v", err)
			}

			return nil
		}()

		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}
