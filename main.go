package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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

func main() {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	appRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()

	apiCnfg := apiCnfg{
		fileserverHits: 0,
	}

	router.Use(middlewareCors)

	setupAppRouter(appRouter, &apiCnfg)
	setupApiRouter(apiRouter, &apiCnfg)
	setupAdminRouter(adminRouter, &apiCnfg)

	router.Mount("/api", apiRouter)
	router.Mount("/app", appRouter)
	router.Mount("/admin", adminRouter)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Failed to start the server")
	}
}
