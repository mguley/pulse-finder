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

	return router
}

// registerHealthCheckRoute defines the health check route.
func registerHealthCheckRoute(router *httprouter.Router, di *application.Container) {
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", di.HealthCheckContainer.Get().HealthCheckHandler.Get().Execute)
}
