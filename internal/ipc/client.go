package ipc

import (
	"encoding/json"
	"fmt"
	"io"
	"jumie/internal/ai"
	"net"
	"os"
	"os/user"
)

type Client struct {
	socketPath string
	Conn       net.Conn
}

type Request struct {
	Type      string   `json:"type"`
	AIMessage string   `json:"ai_message"`
	Commands  []string `json:"commands"`
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

	return &Client{
		socketPath: path,
		Conn:       nil,
	}, nil
}

func (c *Client) openConn() error {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return err
	}

	c.Conn = conn
	return nil
}

func (c *Client) RequestPlan(msg string) (*ai.Plan, error) {
	err := c.openConn()
	if err != nil {
		return nil, err
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			return
		}
	}(c.Conn)

	bytes, err := json.Marshal(Request{"plan", msg, nil})
	if err != nil {
		return nil, err
	}

	_, err = c.Conn.Write(bytes)
	if err != nil {
		return nil, err
	}

	var resp ai.Plan
	err = json.NewDecoder(c.Conn).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DoPlan(plan *ai.Plan) error {
	err := c.openConn()
	if err != nil {
		return err
	}
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			return
		}
	}(c.Conn)

	var commands []string
	for _, step := range plan.Steps {
		commands = append(commands, step.Command)
	}

	req := Request{"exec", "", commands}
	if err := json.NewEncoder(c.Conn).Encode(req); err != nil {
		return err
	}

	_, err = io.Copy(os.Stdout, c.Conn)
	return err
}

func (c *Client) Ping() error {
	err := c.openConn()
	if err != nil {
		return err
	}
	defer c.Conn.Close()

	req := Request{Type: "ping"}
	json.NewEncoder(c.Conn).Encode(req)

	var res struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	json.NewDecoder(c.Conn).Decode(&res)
	if res.Status == "error" {
		return fmt.Errorf(res.Message)
	}
	return nil
}

func (c *Client) StartOllama() error {
	err := c.openConn()
	if err != nil {
		return err
	}
	defer c.Conn.Close()

	req := Request{Type: "start_ollama"}
	json.NewEncoder(c.Conn).Encode(req)

	var res struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	json.NewDecoder(c.Conn).Decode(&res)
	if res.Status == "error" {
		return fmt.Errorf(res.Message)
	}
	return nil
}
