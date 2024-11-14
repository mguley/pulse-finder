package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RouteGroup allows grouping routes and applying middleware to all of them.
type RouteGroup struct {
	router      *httprouter.Router
	middlewares []Middleware
}

// NewRouteGroup creates a new RouteGroup with the given middlewares.
func NewRouteGroup(router *httprouter.Router, middlewares ...Middleware) *RouteGroup {
	return &RouteGroup{router: router, middlewares: middlewares}
}

// HandlerFunc registers a new route with the group, applying the middlewares.
func (rg *RouteGroup) HandlerFunc(method, path string, handlerFunc http.HandlerFunc) {
	h := handlerFunc
	if len(rg.middlewares) > 0 {
		h = ApplyMiddleware(h, rg.middlewares...)
	}
	rg.router.HandlerFunc(method, path, h)
}
