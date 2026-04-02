package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/akaitigo/review-gym/internal/handler"
	"github.com/akaitigo/review-gym/internal/store"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	ms := store.NewMemoryStore()
	h := &handler.Handler{
		Exercises:  ms,
		Reviews:    ms,
		References: ms,
		Scores:     ms,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"status":"ok"}`)
	})

	h.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("review-gym API server starting on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// allowedOrigins returns the set of permitted CORS origins.
func allowedOrigins() map[string]struct{} {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		raw = "http://localhost:3000,http://localhost:5173"
	}
	origins := make(map[string]struct{})
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins[o] = struct{}{}
		}
	}
	return origins
}

// corsMiddleware adds CORS headers using an origin whitelist.
func corsMiddleware(next http.Handler) http.Handler {
	origins := allowedOrigins()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := origins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
