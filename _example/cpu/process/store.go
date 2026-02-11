package process

import (
	"sync"
)

type ProcessStore struct {
	mu   sync.RWMutex
	data []ProcessInfo
}

func (p *ProcessStore) Update(d []ProcessInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data = d
}

func (p *ProcessStore) GetAll() []ProcessInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.data
}
