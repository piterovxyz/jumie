package ai

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"jumie/internal/config"
	"jumie/internal/indexer"
	"net/http"
	"strings"
)

//go:embed ai.md
var system string

type Client struct {
	*http.Client
	cfg                *config.Config
	systemInstructions string
}

type Plan struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
	Format string `json:"format"`
	Stream bool   `json:"stream"`
}
type Response struct {
	Response string `json:"response"`
}

func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		&http.Client{},
		cfg,
		system,
	}, nil
}

func (c *Client) UpdateCache(index indexer.SystemInfo) error {
	paths := strings.Join(index.Path, ", ")
	contextText := fmt.Sprintf(
		"OS: %s\nRelease: %s\nShell: %s\nIs Root: %v\nAvailable Binaries: %s",
		index.OsType, index.OsRelease, index.Shell, index.IsSU, paths,
	)
	c.systemInstructions = cutContext(c.systemInstructions) + "\n" + contextText
	return nil
}

func (c *Client) GeneratePlan(ctx context.Context, query string) (*Plan, error) {
	body := Request{
		Model:  c.cfg.Model,
		Prompt: query,
		System: c.systemInstructions,
		Format: "json",
		Stream: false,
	}

	jbody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(jbody))

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:11434/api/generate", bytes.NewReader(jbody))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	fmt.Println(res.Response)

	result, err := parseResponse(res.Response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func cutContext(prompt string) string {
	marker := "### SYSTEM CONTEXT"
	idx := strings.Index(prompt, marker)
	if idx == -1 {
		return prompt
	}

	return prompt[:idx+len(marker)]
}

func parseResponse(j string) (*Plan, error) {
	cleaned := strings.TrimSpace(j)

	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
	}

	cleaned = strings.TrimSpace(cleaned)

	var plan Plan
	err := json.Unmarshal([]byte(cleaned), &plan)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &plan, nil
}
