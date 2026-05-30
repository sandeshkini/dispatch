package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

const offlineAfter = 90 * time.Second

type Registry struct {
	mu      sync.RWMutex
	workers map[string]*Worker // keyed by worker ID
}

func newRegistry() *Registry {
	return &Registry{workers: make(map[string]*Worker)}
}

// workerID is a stable 8-char hex ID derived from the worker's URL.
func workerID(url string) string {
	h := sha256.Sum256([]byte(url))
	return hex.EncodeToString(h[:4])
}

// Register upserts a worker from a heartbeat payload.
func (r *Registry) Register(reg Registration) *Worker {
	id := workerID(reg.URL)
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	w, exists := r.workers[id]
	if !exists {
		w = &Worker{ID: id, RegisteredAt: now}
		r.workers[id] = w
	}

	w.Label = reg.Label
	w.URL = reg.URL
	w.Version = reg.Version
	w.Capabilities = reg.Capabilities
	w.Sessions = reg.Sessions
	w.LastSeen = now
	w.Online = true

	return w
}

// List returns all workers ordered by label.
func (r *Registry) List() []workerView {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]workerView, 0, len(r.workers))
	for _, w := range r.workers {
		out = append(out, w.view())
	}
	return out
}

// Get returns a single worker by ID.
func (r *Registry) Get(id string) (workerView, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	w, ok := r.workers[id]
	if !ok {
		return workerView{}, false
	}
	return w.view(), true
}

// WorkerURL returns just the base URL for a worker — used by proxy.
func (r *Registry) WorkerURL(id string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	w, ok := r.workers[id]
	if !ok {
		return "", false
	}
	return w.URL, ok
}

// MarkOffline sweeps workers that haven't sent a heartbeat recently.
// Call this from a background goroutine.
func (r *Registry) MarkOffline() {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-offlineAfter)
	for _, w := range r.workers {
		if w.Online && w.LastSeen.Before(cutoff) {
			w.Online = false
		}
	}
}
