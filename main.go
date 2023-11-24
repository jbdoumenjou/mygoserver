package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	apiMetrics := &apiConfig{}

	// Admin routes
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiMetrics.metricsHTMLHandler)

	router.Mount("/admin", adminRouter)

	// Serve static files from the . directory
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	router.Get("/app/assets/", apiMetrics.middlewareMetricsInc(appHandler))
	router.Get("/app", apiMetrics.middlewareMetricsInc(appHandler))

	// API routes
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", healthzHandler)
	apiRouter.Get("/metrics", apiMetrics.metricsHandler)
	apiRouter.Get("/reset", apiMetrics.resetHandler)

	router.Mount("/api", apiRouter)

	corsmux := middlewareCors(router)
	log.Fatal(NewWebServer(":8080", corsmux).Start())
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("serving %s\n", r.Host+r.RequestURI)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: " + fmt.Sprintf("%d", cfg.fileserverHits)))
}

func (cfg *apiConfig) metricsHTMLHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	const html = `<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
`
	w.Write([]byte(fmt.Sprintf(html, cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++

		next.ServeHTTP(w, r)
	}
}
