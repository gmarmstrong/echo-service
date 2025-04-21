package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// responsePayload mirrors the incoming request so callers can verify round‑trip.
type responsePayload struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Host    string              `json:"host"`
	Remote  string              `json:"remote_addr"`
}

// echoHandler writes back basic request details as JSON.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	// TODO Need to handle potential encoding errors?
	_ = enc.Encode(responsePayload{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: r.Header,
		Host:    r.Host,
		Remote:  r.RemoteAddr,
	})
}

// healthzHandler responds with HTTP 204 No Content
func healthzHandler(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusNoContent)
}

func main() {
	port := os.Getenv("PORT")
	// TODO Need to validate port number?
	if port == "" {
		port = "8080"
	}
	// TODO Restrict to localhost or a specific interface/IP address?
	addr := net.JoinHostPort("", port)

	// Set up our ServeMux, an HTTP request multiplexer.
	// It will match the URL of incoming requests against a list
	// of registered patterns and call the handler for the one that
	// most closely matches the URL.
	mux := http.NewServeMux()
	// Register patterns for handlers.
	mux.HandleFunc("/", echoHandler)
	mux.HandleFunc("/healthz", healthzHandler)

	// Use slog (built‑in structured logging).
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 	   5 * time.Second,
		WriteTimeout: 	   10 * time.Second,
		IdleTimeout: 	   60 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Graceful shutdown.
	go func() {
		logger.Info("listening", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	logger.Info("server shut down")
}
