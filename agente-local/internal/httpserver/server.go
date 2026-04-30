package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, allowedOrigins []string, sharedSecret string, handler *Handler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.Healthz)
	mux.HandleFunc("GET /certificates", handler.ListCertificates)
	mux.HandleFunc("POST /sign", handler.Sign)

	wrapped := withCORS(allowedOrigins, withSharedSecret(sharedSecret, mux))
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           wrapped,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func withSharedSecret(sharedSecret string, next http.Handler) http.Handler {
	want := strings.TrimSpace(sharedSecret)
	if want == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions || r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		got := strings.TrimSpace(r.Header.Get("X-Local-Agent-Secret"))
		if got == "" || got != want {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func withCORS(allowedOrigins []string, next http.Handler) http.Handler {
	allowSet := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		v := strings.TrimSpace(origin)
		if v != "" {
			allowSet[v] = struct{}{}
		}
	}
	allowAll := len(allowSet) == 0

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if _, ok := allowSet[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
