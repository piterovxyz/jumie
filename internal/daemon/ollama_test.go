package daemon

import (
	"jumie/internal/installer"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupMockOllama(t *testing.T) {
	t.Setenv("PATH", "")
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	binPath := installer.GetOllamaBinPath()
	err := os.MkdirAll(filepath.Dir(binPath), 0755)
	if err != nil {
		t.Fatal(err)
	}

	mockScript := `#!/bin/sh
if [ "$1" = "serve" ]; then
	sleep 10
elif [ "$1" = "pull" ]; then
	echo "pulling model"
	echo "100%"
fi
`
	err = os.WriteFile(binPath, []byte(mockScript), 0755)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetModelName(t *testing.T) {
	model := GetModelName()
	expected := "gemma4:e2b"
	if runtime.GOOS == "darwin" {
		expected = "gemma4:e2b-mlx"
	}
	if model != expected {
		t.Errorf("expected %q, got %q", expected, model)
	}
}

func TestStartAndStopOllama(t *testing.T) {
	setupMockOllama(t)

	err := StartOllama()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if ollamaCmd == nil {
		t.Fatalf("expected ollamaCmd to be set")
	}

	StopOllama()

	if ollamaCmd != nil {
		t.Fatalf("expected ollamaCmd to be nil after stop")
	}
}

func TestStartOllamaNotInstalled(t *testing.T) {
	t.Setenv("PATH", "")
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	err := StartOllama()
	if err == nil {
		t.Errorf("expected error when not installed, got nil")
	}
}

func TestPullModel(t *testing.T) {
	setupMockOllama(t)

	var progressMessages []string
	progress := func(msg string) {
		progressMessages = append(progressMessages, msg)
	}

	err := PullModel(progress)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(progressMessages) == 0 {
		t.Errorf("expected some progress messages")
	}
}
