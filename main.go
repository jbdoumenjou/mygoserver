package main

import (
	"encoding/json"
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
	apiRouter.Post("/validate_chirp", validateChirpHandler)

	router.Mount("/api", apiRouter)

	corsmux := middlewareCors(router)
	log.Fatal(NewWebServer(":8080", corsmux).Start())
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type resp struct {
		Valid bool   `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Something went wrong: %s", err)
		content, err := json.Marshal(resp{Error: err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(content)

		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		content, err := json.Marshal(resp{Error: "Chirp is too long"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(content)

		return
	}

	w.WriteHeader(http.StatusOK)
	content, err := json.Marshal(resp{Valid: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(content)
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
