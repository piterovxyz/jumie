package indexer

import (
	"runtime"
	"testing"
)

func TestRunIndexer(t *testing.T) {
	cache := NewCache(SystemInfo{})
	RunIndexer(cache)

	info := cache.Get()
	if info.OsType == "" {
		t.Errorf("expected OsType to be set")
	}
}

func TestCheckBinaries(t *testing.T) {
	tools := []string{"ls", "nonexistent-tool-123"}
	res := CheckBinaries(tools)

	if !res["ls"] {
		t.Errorf("expected ls to be found")
	}

	if res["nonexistent-tool-123"] {
		t.Errorf("expected nonexistent-tool-123 to not be found")
	}
}

func TestGetSystemRelease(t *testing.T) {
	osType, osRelease, err := getSystemRelease()
	if err != nil {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			t.Errorf("unexpected error on %s: %v", runtime.GOOS, err)
		}
	} else {
		if osType == "" {
			t.Errorf("expected OsType to be set")
		}
		if osRelease == "" || osRelease == "unsupported" {
			t.Errorf("expected OsRelease to be set and supported")
		}
	}
}
