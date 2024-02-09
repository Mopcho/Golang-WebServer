package main

import (
	"log"
	"net/http"
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

func (apiCnfg *apiCnfg) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCnfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (apiCnfg *apiCnfg) handleMetrictsReset(w http.ResponseWriter, r *http.Request) {
	apiCnfg.fileserverHits = 0

	w.WriteHeader(200)
}

func main() {
	apiCnfg := apiCnfg{
		fileserverHits: 0,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCnfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/", handleHealthz)
	mux.HandleFunc("/metrics", apiCnfg.handleMetricts)
	mux.HandleFunc("/metrics/reset", apiCnfg.handleMetrictsReset)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start the server")
	}
}
