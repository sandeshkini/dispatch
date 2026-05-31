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
	Port      int    `json:"port"`
	AuthToken string `json:"auth_token,omitempty"`
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

	// auth middleware
	var handler http.Handler = mux
	if cfg.AuthToken != "" {
		handler = authMiddleware(cfg.AuthToken, mux)
	}

	// worker registration (called by workers)
	mux.HandleFunc("/api/register", srv.handleRegister)

	// hub API (called by dashboard / external tools)
	mux.HandleFunc("/api/workers", srv.handleWorkers)
	mux.HandleFunc("/api/workers/", func(w http.ResponseWriter, r *http.Request) {
		rest := strings.TrimPrefix(r.URL.Path, "/api/workers/")

		// /api/workers/{id}/ws/{name} → ws info
		if strings.Contains(rest, "/ws/") {
			srv.handleWSInfo(w, r)
			return
		}

		// /api/workers/{id} (no further path) → detail
		if !strings.Contains(rest, "/") {
			srv.handleWorkerDetail(w, r)
			return
		}

		// /api/workers/{id}/{action}[/{name}] → proxy to worker
		srv.handleWorkerProxy(w, r)
	})

	mux.HandleFunc("/session/", srv.handleSessionPage)
	mux.HandleFunc("/health", srv.handleHealth)
	mux.HandleFunc("/", srv.handleDashboard)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("dispatch v%s listening on %s", version, addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /health is always open — used for uptime checks
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		// dashboard and static pages are open when no token configured
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		var bearer string
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			bearer = auth[7:]
		}
		if bearer == "" || bearer != token {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
