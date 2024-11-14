package middleware

import "net/http"

// Middleware defines a type that represents a middleware function.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in a sequence to a handler.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return h
	}

	// Apply middlewares in reverse order to maintain left-to-right chaining
	for i := len(middlewares) - 1; i >= 0; i-- {
		if middlewares[i] != nil {
			h = middlewares[i](h)
		}
	}
	return h
}

// ApplyMiddleware wraps the http.HandlerFunc with the provided middlewares.
func ApplyMiddleware(hf http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var h http.Handler = hf
		h = Chain(h, middlewares...)
		h.ServeHTTP(w, r)
	}
}
