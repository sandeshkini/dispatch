package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(t *testing.T, authToken string) (*server, *httptest.Server) {
	t.Helper()
	reg := newRegistry()
	srv := &server{
		config:   Config{Port: 8888, AuthToken: authToken},
		registry: reg,
	}
	mux := http.NewServeMux()
	var handler http.Handler = mux
	if authToken != "" {
		handler = authMiddleware(authToken, mux)
	}
	mux.HandleFunc("/api/register", srv.handleRegister)
	mux.HandleFunc("/api/workers", srv.handleWorkers)
	mux.HandleFunc("/api/workers/", func(w http.ResponseWriter, r *http.Request) {
		srv.handleWorkerProxy(w, r)
	})
	mux.HandleFunc("/health", srv.handleHealth)
	return srv, httptest.NewServer(handler)
}

func postJSON(t *testing.T, ts *httptest.Server, path string, body any, token string) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", ts.URL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func getJSON(t *testing.T, ts *httptest.Server, path string, token string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest("GET", ts.URL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

// --- Registry tests ---

func TestRegistry_RegisterAndList(t *testing.T) {
	r := newRegistry()
	reg := Registration{
		Label:       "aibo",
		URL:         "http://100.1.1.1:7777",
		WorkerToken: "tok-abc",
		Capabilities: []string{"claude", "terminal"},
		Sessions:    []Session{{Name: "swift-fox", Status: "running", CLI: "claude", Dir: "/home/user"}},
	}
	id := r.Register(reg)
	if id == "" {
		t.Fatal("Register returned empty ID")
	}
	workers := r.List()
	if len(workers) != 1 {
		t.Fatalf("expected 1 worker, got %d", len(workers))
	}
	if workers[0].Label != "aibo" {
		t.Errorf("wrong label: %s", workers[0].Label)
	}
	if !workers[0].Online {
		t.Error("worker should be online")
	}
	if len(workers[0].Sessions) != 1 {
		t.Error("expected 1 session")
	}
}

func TestRegistry_WorkerInfo_Online(t *testing.T) {
	r := newRegistry()
	r.Register(Registration{Label: "test", URL: "http://1.2.3.4:7777", WorkerToken: "secret", Sessions: []Session{}})
	url, token, ok := r.WorkerInfo(workerID("http://1.2.3.4:7777"))
	if !ok {
		t.Fatal("WorkerInfo returned not ok for online worker")
	}
	if url != "http://1.2.3.4:7777" {
		t.Errorf("wrong url: %s", url)
	}
	if token != "secret" {
		t.Errorf("wrong token: %s", token)
	}
}

func TestRegistry_WorkerInfo_Offline(t *testing.T) {
	r := newRegistry()
	id := r.Register(Registration{Label: "test", URL: "http://1.2.3.4:7777", WorkerToken: "tok", Sessions: []Session{}})
	// Force offline
	r.mu.Lock()
	r.workers[id].Online = false
	r.mu.Unlock()

	_, _, ok := r.WorkerInfo(id)
	if ok {
		t.Error("WorkerInfo should return false for offline worker")
	}
}

func TestRegistry_MarkOffline(t *testing.T) {
	r := newRegistry()
	id := r.Register(Registration{Label: "old", URL: "http://10.0.0.1:7777", WorkerToken: "t", Sessions: []Session{}})
	// Backdate LastSeen
	r.mu.Lock()
	r.workers[id].LastSeen = time.Now().Add(-2 * offlineAfter)
	r.mu.Unlock()

	r.MarkOffline()

	w, ok := r.Get(id)
	if !ok {
		t.Fatal("worker should still exist (not evicted yet)")
	}
	if w.Online {
		t.Error("worker should be offline after MarkOffline")
	}
}

func TestRegistry_EvictStale(t *testing.T) {
	r := newRegistry()
	id := r.Register(Registration{Label: "stale", URL: "http://10.0.0.2:7777", WorkerToken: "t", Sessions: []Session{}})
	r.mu.Lock()
	r.workers[id].LastSeen = time.Now().Add(-2 * evictAfter)
	r.mu.Unlock()

	r.MarkOffline()

	_, ok := r.Get(id)
	if ok {
		t.Error("stale worker should have been evicted")
	}
}

func TestRegistry_SliceIsolation(t *testing.T) {
	r := newRegistry()
	sessions := []Session{{Name: "s1", Status: "running", CLI: "claude", Dir: "/"}}
	r.Register(Registration{Label: "w", URL: "http://5.5.5.5:7777", WorkerToken: "t", Sessions: sessions})

	// Mutate original slice — should not affect registry
	sessions[0].Name = "mutated"

	workers := r.List()
	if workers[0].Sessions[0].Name == "mutated" {
		t.Error("registry sessions slice shares backing array with caller — deep copy failed")
	}
}

// --- Auth middleware tests ---

func TestAuth_NoToken_Open(t *testing.T) {
	_, ts := newTestServer(t, "")
	defer ts.Close()
	resp, _ := http.Get(ts.URL + "/health")
	if resp.StatusCode != 200 {
		t.Errorf("expected 200 without auth configured, got %d", resp.StatusCode)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	_, ts := newTestServer(t, "mysecret")
	defer ts.Close()
	resp := postJSON(t, ts, "/api/register", Registration{
		Label: "w", URL: "http://1.1.1.1:7777", WorkerToken: "wt", Sessions: []Session{},
	}, "mysecret")
	if resp.StatusCode != 200 {
		t.Errorf("expected 200 with valid token, got %d", resp.StatusCode)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	_, ts := newTestServer(t, "mysecret")
	defer ts.Close()
	resp := postJSON(t, ts, "/api/register", Registration{
		Label: "w", URL: "http://1.1.1.1:7777", WorkerToken: "wt", Sessions: []Session{},
	}, "wrongtoken")
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 with wrong token, got %d", resp.StatusCode)
	}
}

func TestAuth_MissingToken(t *testing.T) {
	_, ts := newTestServer(t, "mysecret")
	defer ts.Close()
	resp := postJSON(t, ts, "/api/register", Registration{
		Label: "w", URL: "http://1.1.1.1:7777", WorkerToken: "wt", Sessions: []Session{},
	}, "")
	if resp.StatusCode != 401 {
		t.Errorf("expected 401 with missing token, got %d", resp.StatusCode)
	}
}

func TestAuth_HealthAlwaysOpen(t *testing.T) {
	_, ts := newTestServer(t, "supersecret")
	defer ts.Close()
	resp, _ := http.Get(ts.URL + "/health")
	if resp.StatusCode != 200 {
		t.Errorf("health should be unauthenticated, got %d", resp.StatusCode)
	}
}

// --- Register handler tests ---

func TestHandleRegister_Valid(t *testing.T) {
	srv, ts := newTestServer(t, "")
	defer ts.Close()
	postJSON(t, ts, "/api/register", Registration{
		Label: "mac", URL: "http://192.168.0.190:7777", WorkerToken: "tok123",
		Capabilities: []string{"claude"}, Sessions: []Session{},
	}, "")
	workers := srv.registry.List()
	if len(workers) != 1 || workers[0].Label != "mac" {
		t.Errorf("worker not registered correctly: %+v", workers)
	}
}

func TestHandleRegister_MissingURL(t *testing.T) {
	_, ts := newTestServer(t, "")
	defer ts.Close()
	resp := postJSON(t, ts, "/api/register", Registration{Label: "bad"}, "")
	if resp.StatusCode != 400 {
		t.Errorf("expected 400 for missing URL, got %d", resp.StatusCode)
	}
}

func TestHandleRegister_Heartbeat_UpdatesSessions(t *testing.T) {
	srv, ts := newTestServer(t, "")
	defer ts.Close()

	postJSON(t, ts, "/api/register", Registration{
		Label: "w", URL: "http://2.2.2.2:7777", WorkerToken: "t",
		Sessions: []Session{{Name: "s1", Status: "running", CLI: "claude", Dir: "/"}},
	}, "")
	postJSON(t, ts, "/api/register", Registration{
		Label: "w", URL: "http://2.2.2.2:7777", WorkerToken: "t",
		Sessions: []Session{
			{Name: "s1", Status: "stopped", CLI: "claude", Dir: "/"},
			{Name: "s2", Status: "running", CLI: "pi", Dir: "/home"},
		},
	}, "")

	workers := srv.registry.List()
	if len(workers[0].Sessions) != 2 {
		t.Errorf("expected 2 sessions after update, got %d", len(workers[0].Sessions))
	}
}

// --- Workers API tests ---

func TestHandleWorkers_Empty(t *testing.T) {
	_, ts := newTestServer(t, "")
	defer ts.Close()
	resp := getJSON(t, ts, "/api/workers", "")
	var workers []workerView
	json.NewDecoder(resp.Body).Decode(&workers)
	if len(workers) != 0 {
		t.Errorf("expected empty list, got %d workers", len(workers))
	}
}

func TestHandleWorkers_AfterRegister(t *testing.T) {
	_, ts := newTestServer(t, "")
	defer ts.Close()
	postJSON(t, ts, "/api/register", Registration{
		Label: "pi-box", URL: "http://100.5.5.5:7777", WorkerToken: "t",
		Sessions: []Session{},
	}, "")
	resp := getJSON(t, ts, "/api/workers", "")
	var workers []workerView
	json.NewDecoder(resp.Body).Decode(&workers)
	if len(workers) != 1 || workers[0].Label != "pi-box" {
		t.Errorf("unexpected workers: %+v", workers)
	}
}

// --- Health tests ---

func TestHealth(t *testing.T) {
	_, ts := newTestServer(t, "")
	defer ts.Close()
	postJSON(t, ts, "/api/register", Registration{
		Label: "w1", URL: "http://3.3.3.3:7777", WorkerToken: "t", Sessions: []Session{},
	}, "")
	resp := getJSON(t, ts, "/health", "")
	var h map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&h)
	if h["status"] != "ok" {
		t.Errorf("health status not ok: %v", h)
	}
	if h["workers_online"].(float64) != 1 {
		t.Errorf("expected 1 online worker: %v", h)
	}
}
