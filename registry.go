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

	sessions := make([]Session, len(reg.Sessions))
	copy(sessions, reg.Sessions)
	caps := make([]string, len(reg.Capabilities))
	copy(caps, reg.Capabilities)

	w.Label = reg.Label
	w.URL = reg.URL
	w.APIURL = reg.APIURL
	w.Token = reg.WorkerToken
	w.Version = reg.Version
	w.Capabilities = caps
	w.Sessions = sessions
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

// WorkerURL returns the public base URL for an online worker.
func (r *Registry) WorkerURL(id string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.workers[id]
	if !ok || !w.Online {
		return "", false
	}
	return w.URL, true
}

// WorkerWSURL returns the URL to use for browser WebSocket connections.
// Uses APIURL when set (SSO-free register hostname) so the browser can
// connect without needing a Pangolin SSO cookie for the worker's public URL.
// Falls back to the public URL when no APIURL is configured.
func (r *Registry) WorkerWSURL(id string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.workers[id]
	if !ok || !w.Online {
		return "", false
	}
	if w.APIURL != "" {
		return w.APIURL, true
	}
	return w.URL, true
}

// WorkerInfo returns the API URL and auth token for an online worker.
// APIURL is used for hub→worker proxy calls; it falls back to URL when not set.
func (r *Registry) WorkerInfo(id string) (url, token string, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, exists := r.workers[id]
	if !exists || !w.Online {
		return "", "", false
	}
	apiURL := w.APIURL
	if apiURL == "" {
		apiURL = w.URL
	}
	if apiURL == "" {
		return "", "", false
	}
	return apiURL, w.Token, true
}

// Token returns the stored auth token for a worker regardless of online status.
func (r *Registry) Token(id string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if w, ok := r.workers[id]; ok {
		return w.Token
	}
	return ""
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
