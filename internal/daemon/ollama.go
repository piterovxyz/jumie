package daemon

import (
	"fmt"
	"io"
	"jumie/internal/installer"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var ollamaCmd *exec.Cmd

func StartOllama() error {
	binPath := installer.GetOllamaBinPath()
	if !installer.IsOllamaInstalled() {
		return fmt.Errorf("ollama not installed")
	}

	home, _ := os.UserHomeDir()
	modelsDir := filepath.Join(home, ".local", "share", "jumie", "models")

	err := os.MkdirAll(modelsDir, 0755)
	if err != nil {
		return err
	}

	ollamaCmd = exec.Command(binPath, "serve")
	ollamaCmd.Env = os.Environ()
	ollamaCmd.Env = append(ollamaCmd.Env, "OLLAMA_HOST=127.0.0.1:49312")
	ollamaCmd.Env = append(ollamaCmd.Env, "OLLAMA_MODELS="+modelsDir)

	logFile, err := os.OpenFile(filepath.Join(filepath.Dir(modelsDir), "ollama.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		ollamaCmd.Stdout = logFile
		ollamaCmd.Stderr = logFile
	}

	err = ollamaCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start ollama: %w", err)
	}

	log.Printf("[ollama] started isolated instance (pid: %d)", ollamaCmd.Process.Pid)
	return nil
}

func StopOllama() {
	if ollamaCmd != nil && ollamaCmd.Process != nil {
		log.Printf("[ollama] stopping isolated instance (pid: %d)", ollamaCmd.Process.Pid)
		err := ollamaCmd.Process.Kill()
		if err != nil {
			return
		}
		err = ollamaCmd.Wait()
		ollamaCmd = nil
	}
}

func GetModelName() string {
	if runtime.GOOS == "darwin" {
		return "gemma4:e2b-mlx"
	}
	return "gemma4:e2b"
}

func PullModel(progress func(string)) error {
	binPath := installer.GetOllamaBinPath()
	model := GetModelName()

	progress(fmt.Sprintf("pulling model %s...\n", model))

	cmd := exec.Command(binPath, "pull", model)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "OLLAMA_HOST=127.0.0.1:49312")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			progress(string(buf[:n]))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}

	progress(fmt.Sprintf("\nmodel %s downloaded successfully!\n", model))
	return nil
}
