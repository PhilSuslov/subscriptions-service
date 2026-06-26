package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) { w.status = code; w.ResponseWriter.WriteHeader(code) }

type healthLogLimiter struct {
	mu   sync.Mutex
	last time.Time
}

func (l *healthLogLimiter) allow(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.last.IsZero() && now.Sub(l.last) < time.Minute {
		return false
	}
	l.last = now
	return true
}

func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	healthLimiter := &healthLogLimiter{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			if r.Method == http.MethodGet && strings.TrimSuffix(r.URL.Path, "/") == "/health" {
				if !healthLimiter.allow(time.Now()) {
					return
				}
			}

			log.Info("http request", "method", r.Method, "path", r.URL.Path, "status", rw.status, "duration", time.Since(start).String())
		})
	}
}
