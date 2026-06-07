package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetPath(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	path, err := GetPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(tempDir, ".config", "jumie", "config.json")
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	cfg := &Config{Model: "test-model"}
	err := Save(cfg)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Model != cfg.Model {
		t.Errorf("expected model %q, got %q", cfg.Model, loaded.Model)
	}
}

func TestLoadDefault(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("failed to load default config: %v", err)
	}

	expectedModel := "gemma4:e2b"
	if runtime.GOOS == "darwin" {
		expectedModel = "gemma4:e2b-mlx"
	}

	if loaded.Model != expectedModel {
		t.Errorf("expected default model %q, got %q", expectedModel, loaded.Model)
	}

	path, _ := GetPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected default config file to be created, but it was not")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("USERPROFILE", tempDir)

	path, err := GetPath()
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(path, []byte("{invalid json"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Load()
	if err == nil {
		t.Errorf("expected error when loading invalid json, got nil")
	}
}
