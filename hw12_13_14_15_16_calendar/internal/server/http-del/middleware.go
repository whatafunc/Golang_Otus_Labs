package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(rw, r)

			latency := time.Since(start)
			clientIP := r.RemoteAddr
			if ip := r.Header.Get("X-Real-IP"); ip != "" {
				clientIP = ip
			} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
				clientIP = ip
			}
			userAgent := r.UserAgent()
			logLine := fmt.Sprintf(
				"%s - [%s] \"%s %s %s\" %d %v \"%s\"",
				clientIP,
				start.Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.URL.Path,
				r.Proto,
				rw.statusCode,
				latency,
				userAgent,
			)
			logger.Info(logLine)
		})
	}
}
