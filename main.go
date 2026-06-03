package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const version = "0.1.0"

type Config struct {
	Port             int    `json:"port"`
	AuthToken        string `json:"auth_token,omitempty"`
	RegistrationHost string `json:"registration_host,omitempty"` // e.g. "register.dispatch.kingdomofluna.com"
}

type server struct {
	config   Config
	registry *Registry
}

func defaultConfig() Config {
	return Config{Port: 8888}
}

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "dispatch", "config.json")
}

func main() {
	port := flag.Int("port", 0, "port to listen on (overrides config)")
	flag.Parse()

	cfg, err := loadConfig(configPath())
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	if *port != 0 {
		cfg.Port = *port
	}

	reg := newRegistry()

	srv := &server{config: cfg, registry: reg}

	// background: mark workers offline when heartbeats stop
	go func() {
		for range time.Tick(30 * time.Second) {
			reg.MarkOffline()
		}
	}()

	mux := http.NewServeMux()

	// worker registration (called by workers)
	mux.HandleFunc("/api/register", srv.handleRegister)

	// hub API (called by dashboard / external tools)
	mux.HandleFunc("/api/workers", srv.handleWorkers)
	mux.HandleFunc("/api/workers/", func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, "/api/workers/")

		// /api/workers/{id} (no further path) → detail
		if !strings.Contains(rest, "/") {
			srv.handleWorkerDetail(w, r)
			return
		}

		// Split into id + action segment to route correctly.
		// Must check action=="ws" exactly — strings.Contains("/ws/") would
		// misroute paths like /kill/foo/ws/bar to handleWSInfo.
		parts := strings.SplitN(rest, "/", 3)
		if len(parts) >= 2 && parts[1] == "ws" {
			srv.handleWSInfo(w, r)
			return
		}

		// /api/workers/{id}/{action}[/{name}] → proxy to worker
		srv.handleWorkerProxy(w, r)
	})

	mux.HandleFunc("/multi", srv.handleMulti)
	mux.HandleFunc("/session/", srv.handleSessionPage)
	mux.HandleFunc("/health", srv.handleHealth)
	mux.HandleFunc("/", srv.handleDashboard)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("dispatch v%s listening on %s", version, addr)
	if err := http.ListenAndServe(addr, srv.hostMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

// hostMiddleware gates requests based on which hostname they arrive on.
// When RegistrationHost is set, requests from that host are restricted to
// /api/register and /health only — the dashboard and worker API are not
// accessible through the public registration URL.
func (s *server) hostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.RegistrationHost == "" {
			next.ServeHTTP(w, r)
			return
		}
		// Traefik sets X-Forwarded-Host; fall back to Host header.
		host := r.Header.Get("X-Forwarded-Host")
		if host == "" {
			host = r.Host
		}
		// Strip port if present.
		if i := strings.LastIndex(host, ":"); i >= 0 {
			host = host[:i]
		}
		if host == s.config.RegistrationHost {
			if r.URL.Path != "/api/register" && r.URL.Path != "/health" {
				http.NotFound(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
