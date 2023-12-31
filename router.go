package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jbdoumenjou/mygoserver/internal/api/chirp"
	"github.com/jbdoumenjou/mygoserver/internal/api/cors"
	"github.com/jbdoumenjou/mygoserver/internal/api/health"
	"github.com/jbdoumenjou/mygoserver/internal/api/metrics"
	"github.com/jbdoumenjou/mygoserver/internal/api/token"
	"github.com/jbdoumenjou/mygoserver/internal/api/user"
)

type ApiConfig struct {
	JWTSecret string
}

type Storer interface {
	chirp.ChirpStorer
	user.UserStorer
}

func NewRouter(db Storer, tokenManager *token.Manager) http.Handler {
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

	chirpHandler := chirp.NewHandler(db, tokenManager)
	apiRouter.Get("/chirps", chirpHandler.List)
	apiRouter.Get("/chirps/{id}", chirpHandler.Get)
	apiRouter.Delete("/chirps/{id}", chirpHandler.Delete)
	apiRouter.Post("/chirps", chirpHandler.Create)

	userHandler := user.NewHandler(db, tokenManager)
	apiRouter.Post("/users", userHandler.Create)
	apiRouter.Put("/users", userHandler.Update)
	apiRouter.Post("/login", userHandler.Login)
	apiRouter.Post("/refresh", userHandler.Refresh)
	apiRouter.Post("/revoke", userHandler.Revoke)
	apiRouter.Post("/polka/webhooks", userHandler.Upgrade)

	router.Mount("/api", apiRouter)

	return cors.Middleware(router)
}
