package main

import "time"

// ---------------------------------------------------------------------------
// Worker → Hub contract
// Workers POST this to /api/register on startup and every 30s.
// ---------------------------------------------------------------------------

type Registration struct {
	Label        string    `json:"label"`
	URL          string    `json:"url"`
	WorkerToken  string    `json:"worker_token"`
	Version      string    `json:"version,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"` // e.g. ["claude","pi","terminal"]
	Sessions     []Session `json:"sessions"`
}

// ---------------------------------------------------------------------------
// Hub → Worker contract
// Hub calls these endpoints on each worker. Workers must implement them.
//
//   GET  {url}/api/v1/instances
//   POST {url}/api/v1/spawn          body: SpawnRequest
//   POST {url}/api/v1/kill/{name}
//   POST {url}/api/v1/restart/{name}
//   POST {url}/api/v1/resume/{name}
//   GET  {url}/api/v1/output/{name}?lines=100
//   WS   {url}/ws/{name}             (browser connects directly)
// ---------------------------------------------------------------------------

type SpawnRequest struct {
	Dir   string   `json:"dir"`
	CLI   string   `json:"cli,omitempty"`
	Flags []string `json:"flags,omitempty"`
	Name  string   `json:"name,omitempty"`
}

// ---------------------------------------------------------------------------
// Shared types
// ---------------------------------------------------------------------------

type Session struct {
	Name      string `json:"name"`
	Status    string `json:"status"`    // "running" | "stopped"
	CLI       string `json:"cli"`
	Dir       string `json:"dir"`
	Summary   string `json:"summary,omitempty"`
	StartedAt string `json:"started_at,omitempty"`
}

type Worker struct {
	ID           string
	Label        string
	URL          string
	Token        string // worker's own auth token — used by hub when calling back
	Version      string
	Capabilities []string
	Sessions     []Session
	RegisteredAt time.Time
	LastSeen     time.Time
	Online       bool
}

// workerView is the JSON-safe shape served by hub APIs.
type workerView struct {
	ID           string    `json:"id"`
	Label        string    `json:"label"`
	URL          string    `json:"url"`
	Version      string    `json:"version,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"`
	Sessions     []Session `json:"sessions"`
	Online       bool      `json:"online"`
	LastSeen     time.Time `json:"last_seen"`
}

func (w *Worker) view() workerView {
	sessions := w.Sessions
	if sessions == nil {
		sessions = []Session{}
	}
	caps := w.Capabilities
	if caps == nil {
		caps = []string{}
	}
	return workerView{
		ID:           w.ID,
		Label:        w.Label,
		URL:          w.URL,
		Version:      w.Version,
		Capabilities: caps,
		Sessions:     sessions,
		Online:       w.Online,
		LastSeen:     w.LastSeen,
	}
}
