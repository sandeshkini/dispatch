package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var proxyClient = &http.Client{
	Timeout: 15 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// forwardToWorker proxies a request to a worker and writes the response back.
// workerURL is the worker's base URL (e.g. "http://100.x.x.x:7777").
// path is the worker-side path (e.g. "/api/v1/kill/swift-fox").
func forwardToWorker(w http.ResponseWriter, r *http.Request, workerURL, path string) {
	target := strings.TrimRight(workerURL, "/") + path

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, r.Method, target, r.Body)
	if err != nil {
		jsonError(w, "failed to build upstream request: "+err.Error(), 502)
		return
	}
	req.ContentLength = r.ContentLength
	if ct := r.Header.Get("Content-Type"); ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	resp, err := proxyClient.Do(req)
	if err != nil {
		jsonError(w, "worker unreachable: "+err.Error(), 502)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// wsInfoForSession returns the direct WebSocket URL a browser should connect
// to for terminal streaming. The browser connects straight to the worker —
// the hub never relays PTY data.
func wsInfoForSession(workerURL, sessionName string) map[string]string {
	base := strings.TrimRight(workerURL, "/")
	wsBase := strings.Replace(base, "http://", "ws://", 1)
	wsBase = strings.Replace(wsBase, "https://", "wss://", 1)
	return map[string]string{
		"ws_url":     fmt.Sprintf("%s/ws/%s", wsBase, sessionName),
		"worker_url": base,
	}
}

// ---------------------------------------------------------------------------
// JSON helpers
// ---------------------------------------------------------------------------

func jsonOK(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
