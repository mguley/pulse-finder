package route

import (
	"application"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Register initializes all router groups.
func Register(di *application.Container) http.Handler {
	router := httprouter.New()

	registerHealthCheckRoute(router, di)
	registerAuthenticationRoute(router, di)

	return router
}

// registerHealthCheckRoute defines the health check route.
func registerHealthCheckRoute(router *httprouter.Router, di *application.Container) {
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", di.HealthCheckContainer.Get().HealthCheckHandler.Get().Execute)
}

// registerAuthenticationRoute defines the authorization routes.
func registerAuthenticationRoute(router *httprouter.Router, di *application.Container) {
	router.HandlerFunc(http.MethodGet, "/v1/jwt", di.JwtAuthContainer.Get().JwtAuthHandler.Get().Execute)
}
