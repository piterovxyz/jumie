package ipc

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
)

type Client struct {
	socketPath string
	Conn       net.Conn
}

func NewClient() (*Client, error) {
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

	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &Client{
		socketPath: path,
		Conn:       conn,
	}, nil
}

func (c *Client) SendMessage(msg string) error {
	data := Request{
		Payload: json.RawMessage(fmt.Sprintf(`{"ai_message": "%s"}`, msg)),
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(bytes)
	if err != nil {
		return err
	}

	var resp Response
	err = json.NewDecoder(c.Conn).Decode(&resp)
	if err != nil {
		return err
	}

	return nil
}
