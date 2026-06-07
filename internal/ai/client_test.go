package ai

import (
	"context"
	"encoding/json"
	"jumie/internal/config"
	"jumie/internal/indexer"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateCache(t *testing.T) {
	cfg := &config.Config{Model: "test"}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	info := indexer.SystemInfo{
		OsType:    "linux",
		OsRelease: "1.0",
		Shell:     "/bin/bash",
		IsSU:      false,
	}
	tools := map[string]bool{"curl": true, "wget": false}

	err = client.UpdateCache(info, tools)
	if err != nil {
		t.Fatal(err)
	}

	if len(client.systemInstructions) == 0 {
		t.Errorf("expected system instructions to be updated")
	}
}

func TestGeneratePlan(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			Response: "```json\n{\"reasoning\":\"test reasoning\",\"steps\":[{\"command\":\"ls\",\"description\":\"list files\"}]}\n```",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	cfg := &config.Config{Model: "test"}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	client.Client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("tcp", ts.Listener.Addr().String())
		},
	}

	plan, err := client.GeneratePlan(context.Background(), "test query")
	if err != nil {
		t.Fatal(err)
	}

	if plan == nil {
		t.Fatalf("expected plan, got nil")
	}
	if plan.Reasoning != "test reasoning" {
		t.Errorf("expected reasoning 'test reasoning', got %q", plan.Reasoning)
	}
	if len(plan.Steps) != 1 || plan.Steps[0].Command != "ls" {
		t.Errorf("unexpected steps: %+v", plan.Steps)
	}
}

func TestGenerateRecon(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			Response: "{\"tip\":\"use ls\",\"tools\":[\"ls\",\"cat\"]}",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	cfg := &config.Config{Model: "test"}
	client, err := NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	client.Client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("tcp", ts.Listener.Addr().String())
		},
	}

	tools, tip, err := client.GenerateRecon(context.Background(), "test query")
	if err != nil {
		t.Fatal(err)
	}

	if tip != "use ls" {
		t.Errorf("expected tip 'use ls', got %q", tip)
	}
	if len(tools) != 2 || tools[0] != "ls" || tools[1] != "cat" {
		t.Errorf("unexpected tools: %v", tools)
	}
}

func TestParseResponse(t *testing.T) {
	resp1 := "<|think|>reasoning1</|think|>{\"steps\":[]}"
	plan1, err := parseResponse(resp1)
	if err != nil {
		t.Fatal(err)
	}
	if plan1.Reasoning != "reasoning1" {
		t.Errorf("expected reasoning1, got %q", plan1.Reasoning)
	}

	resp2 := "<|channel>thoughtreasoning2<channel|>{\"steps\":[]}"
	plan2, err := parseResponse(resp2)
	if err != nil {
		t.Fatal(err)
	}
	if plan2.Reasoning != "reasoning2" {
		t.Errorf("expected reasoning2, got %q", plan2.Reasoning)
	}
}
