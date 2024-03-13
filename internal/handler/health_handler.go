package handler

import (
	"geolocation-service/docs/swagger"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// HealthHandler is for handler usage in main.go
type HealthHandler struct{}

// NewHealthHandler for kubernetes health and readiness check
func NewHealthHandler(router *chi.Mux) *HealthHandler {
	h := &HealthHandler{}

	router.Get("/swagger/*", httpSwagger.Handler())
	router.Get("/readiness", h.Readiness)
	router.Get("/health", h.Health)
	router.Handle("/metrics", promhttp.Handler())

	swagger.SwaggerInfo.Title = "Geolocation Provider Api"
	swagger.SwaggerInfo.Description = "Microservice for providing service data based on ip address"

	return h
}

// Health godoc
// @Summary Health endpoint for kubernetes health check
// @Tags health-handler
// @Success 200 {} http.Response
// @Router /health [get]
func (*HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readiness godoc
// @Summary Readiness endpoint for kubernetes readiness check
// @Tags health-handler
// @Success 200 {} http.Response
// @Router /readiness [get]
func (*HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
