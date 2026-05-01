package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

var logSkipPrefixes = []string{
	"/_next/",
	"/static/",
	"/themes/",
	"/public/",
}

var logSkipExtensions = []string{
	".js",
	".css",
	".svg",
	".ico",
	".png",
	".woff2",
	".woff",
	".map",
}

type responseLogger struct {
	http.ResponseWriter
	status int
}

func (rw *responseLogger) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Flush repassa ao writer interno quando existir, para SSE/streaming (ex.: /api/ai/chat).
func (rw *responseLogger) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func shouldSkipRequestLog(path string) bool {
	p := strings.ToLower(strings.TrimSpace(path))
	if p == "" {
		return false
	}
	for _, prefix := range logSkipPrefixes {
		if strings.HasPrefix(p, prefix) {
			return true
		}
	}
	for _, ext := range logSkipExtensions {
		if strings.HasSuffix(p, ext) {
			return true
		}
	}
	return false
}

func RequestLoggerFiltered(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldSkipRequestLog(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		lrw := &responseLogger{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, lrw.status, time.Since(start).Round(time.Millisecond))
	})
}

