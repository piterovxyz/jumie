package ipc

import "encoding/json"

type Response struct {
	Command string          `json:"command"`
	Payload json.RawMessage `json:"payload"`
}
