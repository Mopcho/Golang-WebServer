package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type apiCnfg struct {
	fileserverHits int
}

func (apiCnfg *apiCnfg) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiCnfg.fileserverHits++
		next.ServeHTTP(w, r)
	}
}

func (apiCnfg *apiCnfg) handleMetrictsReset(w http.ResponseWriter, r *http.Request) {
	apiCnfg.fileserverHits = 0

	w.WriteHeader(200)
}

func main() {
	router := chi.NewRouter()

	apiCnfg := apiCnfg{
		fileserverHits: 0,
	}
	router.Use(middlewareCors)
	router.Handle("/app", apiCnfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	router.Handle("/app/*", apiCnfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	router.Get("/healthz", handleHealthz)
	router.Get("/metrics", apiCnfg.handleMetricts)
	router.Get("/metrics/reset", apiCnfg.handleMetrictsReset)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start the server")
	}
}
