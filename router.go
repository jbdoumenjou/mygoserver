package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jbdoumenjou/mygoserver/internal/api/chirp"
	"github.com/jbdoumenjou/mygoserver/internal/api/cors"
	"github.com/jbdoumenjou/mygoserver/internal/api/health"
	"github.com/jbdoumenjou/mygoserver/internal/api/metrics"
	"net/http"
)

func NewRouter(chirpStorer chirp.ChirpStorer) http.Handler {
	router := chi.NewRouter()
	apiMetrics := &metrics.Metrics{}

	// Admin routes
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiMetrics.HTMLHandler)

	router.Mount("/admin", adminRouter)

	// Serve static files from the . directory
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	metricsAppHandler := apiMetrics.MiddlewareInc(appHandler)
	router.Get("/app/assets/", metricsAppHandler)
	router.Get("/app", metricsAppHandler)

	// API routes
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", health.Handler)
	apiRouter.Get("/metrics", apiMetrics.TextHandler)
	apiRouter.Get("/reset", apiMetrics.ResetHandler)

	chirpHandler := chirp.NewHandler(chirpStorer)
	apiRouter.Get("/chirps", chirpHandler.List)
	apiRouter.Post("/chirps", chirpHandler.Create)

	router.Mount("/api", apiRouter)

	return cors.Middleware(router)
}
