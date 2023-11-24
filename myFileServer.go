package main

import (
	"log"
	"net/http"
)

func startMyFileServer() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("content")))
	corsmux := middlewareCors(mux)

	log.Fatal(http.ListenAndServe(":8080", corsmux))
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
