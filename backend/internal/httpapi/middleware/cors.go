package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var (
	corsMu      sync.RWMutex
	allowOrigin = []string{}
)

func SetCORSAllowedOrigins(origins []string) {
	clean := make([]string, 0, len(origins))
	for _, o := range origins {
		s := strings.TrimSpace(o)
		if s == "" {
			continue
		}
		clean = append(clean, s)
	}
	corsMu.Lock()
	allowOrigin = clean
	corsMu.Unlock()
}

// CORS enables browser requests from the local frontend in development.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			if !isAllowedOrigin(origin) {
				if r.Method == http.MethodOptions {
					http.Error(w, "origin nao permitido", http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return true
	}
	if isLocalOrWailsOrigin(origin) {
		return true
	}
	corsMu.RLock()
	defer corsMu.RUnlock()
	for _, allowed := range allowOrigin {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}

func isLocalOrWailsOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if strings.EqualFold(u.Scheme, "wails") {
		return true
	}

	host := strings.ToLower(strings.TrimSpace(u.Host))
	if host == "" {
		return false
	}
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") || strings.HasPrefix(host, "[::1]") {
		return true
	}
	h, _, err := net.SplitHostPort(host)
	if err == nil {
		h = strings.Trim(strings.ToLower(h), "[]")
		return h == "localhost" || h == "127.0.0.1" || h == "::1"
	}
	return false
}
