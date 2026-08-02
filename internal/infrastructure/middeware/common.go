package middeware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for counter := len(middlewares) - 1; counter >= 0; counter-- {
		handler = middlewares[counter](handler)
	}
	return handler
}
