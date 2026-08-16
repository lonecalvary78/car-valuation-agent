package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter

	statusCode  int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status
		}

		next.ServeHTTP(wrapped, r)

		//nolint:gosec // method/URI are sanitized by sanitizeForLog to strip CR/LF before logging
		log.Printf("%d %q %q %v",
			wrapped.statusCode,
			sanitizeForLog(r.Method),
			sanitizeForLog(r.RequestURI),
			time.Since(start),
		)
	})
}
