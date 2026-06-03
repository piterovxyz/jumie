package ipc

import "encoding/json"

type Request struct {
	Payload json.RawMessage `json:"payload"`
}
