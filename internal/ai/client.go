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

//go:embed recon.md
var reconSystem string

const ollamaUrl = "http://127.0.0.1:11434/api/generate"

type Client struct {
	*http.Client
	cfg                *config.Config
	systemInstructions string
}

type Plan struct {
	Reasoning string `json:"reasoning"`
	Steps     []Step `json:"steps"`
}

type Step struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Request struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	System  string         `json:"system"`
	Format  string         `json:"format"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options,omitempty"`
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

func (c *Client) UpdateCache(index indexer.SystemInfo, checkedTools map[string]bool) error {
	var toolsInfo []string
	for tool, exists := range checkedTools {
		status := "missing"
		if exists {
			status = "installed"
		}
		toolsInfo = append(toolsInfo, fmt.Sprintf("- %s: %s", tool, status))
	}

	contextText := fmt.Sprintf(
		"OS: %s\nRelease: %s\nShell: %s\nIs Root: %v\n\nChecked Tools:\n%s",
		index.OsType, index.OsRelease, index.Shell, index.IsSU, strings.Join(toolsInfo, "\n"),
	)

	c.systemInstructions = cutContext(c.systemInstructions) + "\n### SYSTEM CONTEXT\n" + contextText
	return nil
}

func cutContext(prompt string) string {
	marker := "### SYSTEM CONTEXT"
	idx := strings.Index(prompt, marker)
	if idx == -1 {
		return prompt
	}

	return prompt[:idx+len(marker)]
}

func (c *Client) GeneratePlan(ctx context.Context, query string) (*Plan, error) {
	body := Request{
		Model:   c.cfg.Model,
		Prompt:  query,
		System:  c.systemInstructions,
		Stream:  false,
		Options: map[string]any{"num_ctx": 8192},
	}

	jbody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaUrl, bytes.NewReader(jbody))
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

	result, err := parseResponse(res.Response)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GenerateRecon(ctx context.Context, query string) ([]string, error) {
	body := Request{
		Model:   c.cfg.Model,
		Prompt:  query,
		System:  reconSystem,
		Stream:  false,
		Format:  "json",
		Options: map[string]any{"num_ctx": 8192},
	}

	jbody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaUrl, bytes.NewReader(jbody))
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
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	j := res.Response

	j = strings.TrimPrefix(strings.TrimSpace(j), "```json")
	j = strings.TrimPrefix(strings.TrimSpace(j), "```")
	j = strings.TrimSuffix(strings.TrimSpace(j), "```")

	first := strings.Index(j, "{")
	last := strings.LastIndex(j, "}")
	if first != -1 && last != -1 && last >= first {
		j = j[first : last+1]
	}
	var reconRes struct {
		Tools []string `json:"tools"`
	}

	if err := json.Unmarshal([]byte(strings.TrimSpace(j)), &reconRes); err != nil {
		return nil, fmt.Errorf("failed to parse recon tools: %w\nRaw string: %s", err, j)
	}
	return reconRes.Tools, nil
}

func parseResponse(j string) (*Plan, error) {
	var reasoning string

	if start := strings.Index(j, "<|channel>thought"); start != -1 {
		if end := strings.Index(j, "<channel|>"); end != -1 && end > start {
			reasoning = strings.TrimSpace(j[start+17 : end])
			j = strings.TrimSpace(j[end+10:])
		}
	}

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

	if reasoning != "" {
		plan.Reasoning = reasoning
	}

	return &plan, nil
}
