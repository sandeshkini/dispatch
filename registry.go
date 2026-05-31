package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

const offlineAfter = 90 * time.Second
const evictAfter = 1 * time.Hour

type Registry struct {
	mu      sync.RWMutex
	workers map[string]*Worker
}

func newRegistry() *Registry {
	return &Registry{workers: make(map[string]*Worker)}
}

func workerID(url string) string {
	h := sha256.Sum256([]byte(url))
	return hex.EncodeToString(h[:4])
}

// Register upserts a worker from a heartbeat payload.
func (r *Registry) Register(reg Registration) string {
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
	w.Token = reg.WorkerToken
	w.Version = reg.Version
	w.Capabilities = reg.Capabilities
	w.Sessions = reg.Sessions
	w.LastSeen = now
	w.Online = true

	return id
}

// List returns a snapshot of all workers.
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

// WorkerURL returns the base URL for an online worker.
func (r *Registry) WorkerURL(id string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.workers[id]
	if !ok || !w.Online {
		return "", false
	}
	return w.URL, true
}

// WorkerInfo returns the URL and auth token for an online worker.
func (r *Registry) WorkerInfo(id string) (url, token string, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, exists := r.workers[id]
	if !exists || !w.Online {
		return "", "", false
	}
	return w.URL, w.Token, true
}

// MarkOffline marks stale workers offline and evicts workers unseen for over an hour.
func (r *Registry) MarkOffline() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	offlineCutoff := now.Add(-offlineAfter)
	evictCutoff := now.Add(-evictAfter)

	for id, w := range r.workers {
		if w.LastSeen.Before(evictCutoff) {
			delete(r.workers, id)
			continue
		}
		if w.Online && w.LastSeen.Before(offlineCutoff) {
			w.Online = false
		}
	}
}
