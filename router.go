package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jbdoumenjou/mygoserver/handler"
	"net/http"
)

func NewRouter() http.Handler {
	router := chi.NewRouter()
	apiMetrics := &handler.Metrics{}

	// Admin routes
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiMetrics.MetricsHTMLHandler)

	router.Mount("/admin", adminRouter)

	// Serve static files from the . directory
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	metricsAppHandler := apiMetrics.MiddlewareInc(appHandler)
	router.Get("/app/assets/", metricsAppHandler)
	router.Get("/app", metricsAppHandler)

	// API routes
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handler.Healthz)
	apiRouter.Get("/metrics", apiMetrics.MetricsTextHandler)
	apiRouter.Get("/reset", apiMetrics.ResetHandler)
	apiRouter.Post("/validate_chirp", handler.ValidateChirp)

	router.Mount("/api", apiRouter)

	return handler.CORSMiddleware(router)
}
