package ipc

import (
	"encoding/json"
	"net"
	"os"
	"os/user"
)

type Client struct {
	socketPath string
	Conn       net.Conn
}

type msgPayload struct {
	AIMessage string `json:"ai_message"`
}

func NewClient() (*Client, error) {
	var path string

	current, err := user.Current()
	if err != nil {
		return nil, err
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

func (c *Client) SendMessage(msg string) (*Response, error) {
	bytes, err := json.Marshal(msgPayload{msg})
	if err != nil {
		return nil, err
	}

	data := Request{
		Payload: json.RawMessage(bytes),
	}

	bytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}

	_, err = c.Conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.NewDecoder(c.Conn).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
