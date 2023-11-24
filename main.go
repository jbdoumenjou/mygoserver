package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", myHandlerfunc)
	corsMux := middlewareCors(mux)

	srv := http.Server{Addr: ":8080", Handler: corsMux}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}

	log.Println("server closed")
}

func myHandlerfunc(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
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
