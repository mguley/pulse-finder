package route

import (
	"application"
	"interfaces/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Register initializes all router groups.
func Register(di *application.Container) http.Handler {
	router := httprouter.New()

	// Register unprotected route for JWT generation
	registerAuthenticationRoute(router, di)

	middlewares := []middleware.Middleware{
		di.InterfacesContainer.Get().JwtAuthMiddleware.Get().Handle,
	}
	// Create a RouteGroup for protected routes
	protectedGroup := middleware.NewRouteGroup(router, middlewares...)

	// Register protected routes
	registerHealthCheckRoute(protectedGroup, di)
	registerVacancyRoutes(protectedGroup, di)

	return router
}

// registerHealthCheckRoute defines the health check route.
func registerHealthCheckRoute(rg *middleware.RouteGroup, di *application.Container) {
	rg.HandlerFunc(http.MethodGet, "/v1/healthcheck", di.HealthCheckContainer.Get().HealthCheckHandler.Get().Execute)
}

// registerAuthenticationRoute defines the route for generating JWT tokens.
func registerAuthenticationRoute(router *httprouter.Router, di *application.Container) {
	router.HandlerFunc(http.MethodGet, "/v1/jwt", di.JwtAuthContainer.Get().JwtAuthHandler.Get().Execute)
}

// registerVacancyRoutes defines vacancy related routes.
func registerVacancyRoutes(rg *middleware.RouteGroup, di *application.Container) {
	const (
		vacancyCreate = "/v1/vacancies"
		vacancyGet    = "/v1/vacancies/:id"
		vacancyDelete = "/v1/vacancies/:id"
		vacancyPatch  = "/v1/vacancies/:id"
	)
	rg.HandlerFunc(http.MethodPost, vacancyCreate, di.VacancyContainer.Get().CreateHandler.Get().Execute)
	rg.HandlerFunc(http.MethodGet, vacancyGet, di.VacancyContainer.Get().GetHandler.Get().Execute)
	rg.HandlerFunc(http.MethodDelete, vacancyDelete, di.VacancyContainer.Get().DeleteHandler.Get().Execute)
	rg.HandlerFunc(http.MethodPatch, vacancyPatch, di.VacancyContainer.Get().UpdateHandler.Get().Execute)
}
