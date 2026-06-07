package indexer

import (
	"testing"
)

func TestInfoCache(t *testing.T) {
	initial := SystemInfo{
		OsType:    "linux",
		OsRelease: "1.0",
		Shell:     "/bin/sh",
		IsSU:      false,
	}

	cache := NewCache(initial)

	got := cache.Get()
	if got.OsType != initial.OsType {
		t.Errorf("expected %q, got %q", initial.OsType, got.OsType)
	}

	updated := SystemInfo{
		OsType:    "darwin",
		OsRelease: "14.0",
		Shell:     "/bin/zsh",
		IsSU:      true,
	}

	cache.Write(updated)

	got = cache.Get()
	if got.OsType != updated.OsType {
		t.Errorf("expected %q, got %q", updated.OsType, got.OsType)
	}
}
