package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthzHandler)

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	apiMetrics := &apiConfig{}
	mux.Handle("/app/", apiMetrics.middlewareMetricsInc(appHandler))
	mux.HandleFunc("/metrics", apiMetrics.metricsHandler)
	mux.HandleFunc("/reset", apiMetrics.resetHandler)

	corsmux := middlewareCors(mux)

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

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++

		next.ServeHTTP(w, r)
	})
}
