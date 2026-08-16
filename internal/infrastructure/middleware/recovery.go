package middleware

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}
			recoveredErr, ok := recovered.(error)
			if ok && errors.Is(recoveredErr, http.ErrAbortHandler) {
				panic(recovered)
			}
			//nolint:gosec // method/URI are sanitized by sanitizeForLog to strip CR/LF before logging
			log.Printf("panic recovered [method: %q, uri: %q, panic: %v]\n%s",
				sanitizeForLog(r.Method), sanitizeForLog(r.RequestURI), recovered, debug.Stack())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}
