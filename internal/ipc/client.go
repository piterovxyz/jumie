package ipc

import (
	"net"
)

type Client struct {
	socketPath string
	Conn       net.Conn
}

func NewClient(socketPath string) (*Client, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}

	return &Client{
		socketPath: socketPath,
		Conn:       conn,
	}, nil
}
