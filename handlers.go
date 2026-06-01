package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "POST only", 405)
		return
	}
	if s.config.AuthToken != "" {
		bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if bearer != s.config.AuthToken {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", 401)
			return
		}
	}
	var reg Registration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		jsonError(w, "invalid JSON: "+err.Error(), 400)
		return
	}
	if reg.URL == "" {
		jsonError(w, "url required", 400)
		return
	}
	id := s.registry.Register(reg)
	jsonOK(w, map[string]string{"id": id, "status": "registered"})
}

func (s *server) handleWorkers(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, s.registry.List())
}

func (s *server) handleWorkerDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/workers/")
	id = strings.SplitN(id, "/", 2)[0]
	worker, ok := s.registry.Get(id)
	if !ok {
		jsonError(w, "worker not found", 404)
		return
	}
	jsonOK(w, worker)
}

// handleWorkerProxy forwards worker actions: spawn, kill, restart, resume, output.
// URL pattern: /api/workers/{id}/{action}[/{name}]
func (s *server) handleWorkerProxy(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/workers/")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) < 2 {
		jsonError(w, "invalid path", 400)
		return
	}
	id, tail := parts[0], parts[1]

	workerURL, workerToken, ok := s.registry.WorkerInfo(id)
	if !ok {
		jsonError(w, "worker not found", 404)
		return
	}

	workerPath := "/api/v1/" + tail
	if r.URL.RawQuery != "" {
		workerPath += "?" + r.URL.RawQuery
	}

	forwardToWorker(w, r, workerURL, workerToken, workerPath)
}

// handleWSInfo returns the direct WebSocket URL for a session.
// Browser uses this to connect to the worker PTY without going through the hub.
// GET /api/workers/{id}/ws/{name}
func (s *server) handleWSInfo(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/workers/")
	// expects: {id}/ws/{name}
	parts := strings.SplitN(rest, "/ws/", 2)
	if len(parts) != 2 {
		jsonError(w, "invalid path — expected /api/workers/{id}/ws/{name}", 400)
		return
	}
	id := strings.SplitN(parts[0], "/", 2)[0]
	name := parts[1]

	workerURL, ok := s.registry.WorkerURL(id)
	if !ok {
		jsonError(w, "worker not found", 404)
		return
	}

	jsonOK(w, wsInfoForSession(workerURL, name))
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	workers := s.registry.List()
	online := 0
	for _, wk := range workers {
		if wk.Online {
			online++
		}
	}
	jsonOK(w, map[string]any{
		"status":          "ok",
		"workers_total":   len(workers),
		"workers_online":  online,
		"workers_offline": len(workers) - online,
	})
}

func (s *server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	dashTmpl.Execute(w, nil)
}

func (s *server) handleSessionPage(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/session/")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		http.NotFound(w, r)
		return
	}
	workerID, sessionName := parts[0], parts[1]

	worker, ok := s.registry.Get(workerID)
	if !ok {
		http.Error(w, "worker not found", 404)
		return
	}

	status := "stopped"
	for _, sess := range worker.Sessions {
		if sess.Name == sessionName {
			status = sess.Status
			break
		}
	}

	info := wsInfoForSession(worker.URL, sessionName)
	token := s.registry.Token(workerID)
	data := newSessionData(workerID, worker.Label, sessionName, status, info["ws_url"], token)
	sessionTmpl.Execute(w, data)
}
