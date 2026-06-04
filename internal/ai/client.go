package ai

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"jumie/internal/indexer"
	"strings"

	"google.golang.org/genai"
)

//go:embed ai.md
var system string

type Client struct {
	*genai.Client
	model              string
	apiKey             string
	systemInstructions string
}

type Plan struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("api key cannot be empty")
	}

	model := "gemini-3.1-flash-lite"

	ctx := context.Background()
	gc, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		gc,
		model,
		apiKey,
		system,
	}, nil
}

func (c *Client) ValidateKey(ctx context.Context) error {
	_, err := c.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text("ping"),
		nil,
	)
	if err != nil {
		return fmt.Errorf("api key validation failed: %w", err)
	}
	return nil
}

func (c *Client) UpdateCache(index indexer.SystemInfo) error {
	info, err := json.Marshal(index)
	if err != nil {
		return err
	}

	c.systemInstructions = cutContext(c.systemInstructions) + "\n" + string(info)
	return nil
}

func (c *Client) GeneratePlan(ctx context.Context, query string) (*Plan, error) {
	response, err := c.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(query),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{
					{Text: c.systemInstructions},
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		res, err := parseResponse(response.Candidates[0].Content.Parts[0].Text)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		return res, nil
	}

	return nil, errors.New("received empty response from ai model")
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
