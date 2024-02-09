package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setupAppRouter(appRouter *chi.Mux, apiCnfg *apiCnfg) {
	const filepathRoot = "."

	fsHandler := apiCnfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	appRouter.Handle("/", fsHandler)
	appRouter.Handle("/*", fsHandler)
}
