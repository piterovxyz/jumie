package indexer

import "sync"

type SystemInfo struct {
	OsType    string
	OsRelease string
	Path      []string
	Meta      map[string]string
}

type InfoCache struct {
	mu   sync.RWMutex
	info SystemInfo
}

func NewCache(info SystemInfo) *InfoCache {
	return &InfoCache{
		mu:   sync.RWMutex{},
		info: info,
	}
}

func (i *InfoCache) Get() SystemInfo {
	i.mu.RLock()
	defer i.mu.RUnlock()

	return i.info
}

func (i *InfoCache) Write(sys SystemInfo) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.info = sys
}
