package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, mw := range slices.Backward(middlewares) {
		handler = mw(handler)
	}
	return handler
}

// Combine folds several middlewares into one, applied in the given order
// (the first middleware runs outermost, i.e. first).
func Combine(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		return Chain(next, middlewares...)
	}
}
