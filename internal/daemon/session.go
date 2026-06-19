package daemon

import "sync"

type SessionRegistry interface {
	Register(sessionID string, port int) error
	Resolve(port int) (sessionID string, ok bool)
	Release(port int) error
}

type MemorySessionRegistry struct {
	mu     sync.RWMutex
	byPort map[int]string
	byID   map[string]int
	active int
}

func NewMemorySessionRegistry() *MemorySessionRegistry {
	return &MemorySessionRegistry{
		byPort: map[int]string{},
		byID:   map[string]int{},
	}
}

func (r *MemorySessionRegistry) Register(sessionID string, port int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byPort[port] = sessionID
	r.byID[sessionID] = port
	r.active = len(r.byID)
	return nil
}

func (r *MemorySessionRegistry) Resolve(port int) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sessionID, ok := r.byPort[port]
	return sessionID, ok
}

func (r *MemorySessionRegistry) Release(port int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if sessionID, ok := r.byPort[port]; ok {
		delete(r.byPort, port)
		delete(r.byID, sessionID)
	}
	r.active = len(r.byID)
	return nil
}

func (r *MemorySessionRegistry) Active() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.active
}
