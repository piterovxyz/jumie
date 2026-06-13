package installer

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetOllamaBinPath(t *testing.T) {
	t.Setenv("PATH", "")
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	var expected string
	if runtime.GOOS == "darwin" {
		expected = filepath.Join(tempDir, ".local", "share", "jumie", "ollama")
	} else {
		expected = filepath.Join(tempDir, ".local", "share", "jumie", "bin", "ollama")
	}

	actual := GetOllamaBinPath()

	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestIsOllamaInstalled(t *testing.T) {
	t.Setenv("PATH", "")
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	if IsOllamaInstalled() {
		t.Errorf("expected false, got true")
	}

	binPath := GetOllamaBinPath()
	err := os.MkdirAll(filepath.Dir(binPath), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(binPath, []byte("dummy"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	if !IsOllamaInstalled() {
		t.Errorf("expected true, got false")
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestInstallOllama(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	mockBinDir := filepath.Join(tempDir, "mockbin")
	err := os.MkdirAll(mockBinDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	mockTarPath := filepath.Join(mockBinDir, "tar")
	mockTarContent := `#!/bin/sh
		exit 0
	`

	err = os.WriteFile(mockTarPath, []byte(mockTarContent), 0755)
	if err != nil {
		t.Fatal(err)
	}

	origPath := os.Getenv("PATH")
	t.Setenv("PATH", mockBinDir+string(os.PathListSeparator)+origPath)

	var progressMessages []string
	progress := func(msg string) {
		progressMessages = append(progressMessages, msg)
	}

	mockClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		}),
	}

	err = InstallOllama(context.Background(), mockClient, progress)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(progressMessages) != 2 {
		t.Errorf("expected 2 progress messages, got %d", len(progressMessages))
	}
}
