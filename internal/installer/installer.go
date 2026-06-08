package installer

import (
	"fmt"
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

func InstallOllama(progress func(string)) error {
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

	progress(fmt.Sprintf("downloading ollama archive for %s...", runtime.GOOS))

	cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -fL --retry 3 --retry-delay 2 %s | tar --zstd -xf - -C %s", url, targetDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download or extraction failed: %w", err)
	}

	progress("ollama downloaded and extracted successfully!")
	return nil
}
