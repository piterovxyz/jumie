package installer

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func GetOllamaBinPath() string {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, ".local", "share", "jumie", "ollama")
	}
	return filepath.Join(home, ".local", "share", "jumie", "bin", "ollama")
}

func IsOllamaInstalled() bool {
	_, err := os.Stat(GetOllamaBinPath())
	return err == nil
}

func InstallOllama(ctx context.Context, client *http.Client, progress func(string)) error {
	home, _ := os.UserHomeDir()
	targetDir := filepath.Join(home, ".local", "share", "jumie")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return err
	}

	url := "https://github.com/ollama/ollama/releases/latest/download/ollama-linux-amd64.tar.zst"
	if runtime.GOOS == "darwin" {
		url = "https://github.com/ollama/ollama/releases/latest/download/ollama-darwin.tgz"
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "arm64" {
		url = "https://github.com/ollama/ollama/releases/latest/download/ollama-linux-arm64.tar.zst"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	progress(fmt.Sprintf("downloading ollama archive for %s...", runtime.GOOS))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var tarArgs []string
	if runtime.GOOS == "darwin" {
		tarArgs = []string{"-zxf", "-", "-C", targetDir}
	} else {
		tarArgs = []string{"--zstd", "-xf", "-", "-C", targetDir}
	}

	cmd := exec.CommandContext(ctx, "tar", tarArgs...)
	cmd.Stdin = resp.Body
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download or extraction failed: %w", err)
	}

	progress("ollama downloaded and extracted successfully!")
	return nil
}
