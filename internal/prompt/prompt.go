package prompt

import (
	_ "embed"
	"encoding/json"
	"jumie/internal/indexer"
)

//go:embed ai.md
var systemPrompt string

type Message struct {
	SystemPrompt string `json:"system_prompt"`
	UserPrompt   string `json:"user_prompt"`
}

func (m *Message) Build(msg string, index *indexer.SystemInfo) error {
	info, err := json.Marshal(index)
	if err != nil {
		return nil
	}

	systemPrompt = systemPrompt + string(info)

	m.SystemPrompt = systemPrompt
	m.UserPrompt = msg

	return nil
}
